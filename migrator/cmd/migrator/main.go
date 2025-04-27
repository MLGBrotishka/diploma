package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"migrator/config"
	grpc_client "migrator/internal/adapter/grpc"
	"migrator/internal/adapter/repository/intiter"
	"migrator/internal/adapter/repository/migration"
	"migrator/internal/service/initializer"
	migratorService "migrator/internal/service/migrator"
	"migrator/pkg/api/migrator"
	"migrator/pkg/logger"
	"migrator/pkg/postgres"

	"github.com/rs/cors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	defaultConfigPath = "config/config.yaml"
)

func main() {
	// Получение пути к файлу конфигурации
	path := flag.String("config", defaultConfigPath, "path to config file")

	flag.Parse()

	cfg, err := config.NewConfig(*path)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	dbConn, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer dbConn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	initerRepo := intiter.New(dbConn.Pool)
	initerSrv := initializer.New(initerRepo)
	err = initerSrv.InitDB(ctx)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	migrationRepo := migration.New(dbConn.Pool)
	migrationSrv := migratorService.New(migrationRepo)

	grpcService := grpc_client.NewMigration(migrationSrv)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("Recovered from panic", p)

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger.New(logger.InfoLevel)), loggingOpts...),
	))

	reflection.Register(grpcServer)

	migrator.RegisterMigrationServiceServer(grpcServer, grpcService)

	mux := runtime.NewServeMux()
	err = migrator.RegisterMigrationServiceHandlerServer(ctx, mux, grpcService)
	if err != nil {
		log.Fatalf("failed to register handler: %v", err)
	}

	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: withCors,
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	httpServer.ListenAndServe()

	grpcErrChan := make(chan error)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			grpcErrChan <- err
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	logger.Info("Server started, waiting for requests...")
	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-grpcErrChan:
		logger.Error(fmt.Errorf("app - Run - grpcServer.Serve: %w", err))
	}

	// Shutdown
	grpcServer.GracefulStop()
}

// InterceptorLogger adapts logger to interceptor logger.
func InterceptorLogger(l *logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Debug(fmt.Sprintf("%s: %s", lvl, msg), fields...)
	})
}
