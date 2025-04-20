package grpc_client

import (
	"context"

	"migrator/pkg/api/migrator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MigrationService interface{}

type Service struct {
	migrator.UnimplementedMigrationServiceServer
	srv MigrationService
}

func NewMigration(srv MigrationService) *Service {
	return &Service{
		srv: srv,
	}
}

func (s *Service) ApplyMigration(ctx context.Context, req *migrator.ApplyMigrationRequest) (*migrator.ApplyMigrationResponse, error) {
	migrationID := req.GetMigrationId()
	userID := req.GetUserId()

	// Проверяем правильность идентификаторов
	if migrationID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "migration_id must be greater than 0")
	}
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0")
	}

	err := s.srv.ApplyMigration(ctx, migrationID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &migrator.CreateMigrationResponse{}, nil
}

func (s *Service) CreateMigration(ctx context.Context, req *migrator.CreateMigrationRequest) (*migrator.CreateMigrationResponse, error) {
}

func (s *Service) GetMigrationStatus(ctx context.Context, req *migrator.GetMigrationStatusRequest) (*migrator.GetMigrationStatusResponse, error) {
}

func (s *Service) ListMigrations(ctx context.Context, req *migrator.ListMigrationsRequest) (*migrator.ListMigrationsResponse, error) {
}

func (s *Service) RollbackMigration(ctx context.Context, req *migrator.RollbackMigrationRequest) (*migrator.RollbackMigrationResponse, error) {
}
