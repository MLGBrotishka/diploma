package grpc_client

import (
	"context"
	"time"

	"migrator/pkg/api/migrator"

	"migrator/internal/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MigrationService interface {
	CreateMigration(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error)
	ApplyMigration(ctx context.Context, migrationIDs []int64, userID int64) (time.Time, error)
	RollbackMigration(ctx context.Context, migrationID, userID int64) (time.Time, error)
	ListMigrations(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error)
	GetMigration(ctx context.Context, migrationID int64) (entity.MigrationInfo, error)
}

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
	migrationIDs := req.GetMigrationIds()
	userID := req.GetUserId()

	// Проверяем правильность идентификаторов
	if len(migrationIDs) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "migration_ids must be greater than 0")
	}
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0")
	}

	appliedAt, err := s.srv.ApplyMigration(ctx, migrationIDs, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &migrator.ApplyMigrationResponse{AppliedAt: appliedAt.Format(time.DateTime)}, nil
}

func (s *Service) CreateMigration(ctx context.Context, req *migrator.CreateMigrationRequest) (*migrator.CreateMigrationResponse, error) {
	name := req.GetName()
	description := req.GetDescription()
	script := req.GetScript()
	rollbackScript := req.GetRollbackScript()
	userID := req.GetUserId()

	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name cannot be empty")
	}
	if script == "" {
		return nil, status.Errorf(codes.InvalidArgument, "script cannot be empty")
	}
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0")
	}

	migrationID, err := s.srv.CreateMigration(ctx, name, description, script, rollbackScript, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &migrator.CreateMigrationResponse{MigrationId: migrationID}, nil
}

func (s *Service) GetMigration(ctx context.Context, req *migrator.GetMigrationRequest) (*migrator.GetMigrationResponse, error) {
	migrationID := req.GetMigrationId()

	if migrationID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "migration_id must be greater than 0")
	}

	migration, err := s.srv.GetMigration(ctx, migrationID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &migrator.GetMigrationResponse{Migration: convertToGrpcMigration(migration)}, nil
}

func (s *Service) ListMigrations(ctx context.Context, req *migrator.ListMigrationsRequest) (*migrator.ListMigrationsResponse, error) {
	statusFilter := req.GetStatus()

	migrations, err := s.srv.ListMigrations(ctx, statusFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	result := convertToGrpcMigrations(migrations)

	return &migrator.ListMigrationsResponse{Migrations: result}, nil
}

func convertToGrpcMigrations(migrations []entity.MigrationInfo) []*migrator.MigrationInfo {
	result := make([]*migrator.MigrationInfo, len(migrations))
	for i, migration := range migrations {
		result[i] = convertToGrpcMigration(migration)
	}
	return result
}

func convertToGrpcMigration(migration entity.MigrationInfo) *migrator.MigrationInfo {
	return &migrator.MigrationInfo{
		Id:              migration.ID,
		Name:            migration.Name,
		Description:     migration.Description,
		Script:          migration.Script,
		RollbackScript:  migration.RollbackScript,
		Status:          migration.Status.String(),
		CreatedBy:       migration.CreatedBy,
		StatusUpdatedAt: migration.StatusUpdatedAt.Format(time.DateTime),
	}
}

func (s *Service) RollbackMigration(ctx context.Context, req *migrator.RollbackMigrationRequest) (*migrator.RollbackMigrationResponse, error) {
	migrationID := req.GetMigrationId()
	userID := req.GetUserId()

	if migrationID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "migration_id must be greater than 0")
	}
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0")
	}

	rolledBackAt, err := s.srv.RollbackMigration(ctx, migrationID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &migrator.RollbackMigrationResponse{RolledBackAt: rolledBackAt.Format(time.DateTime)}, nil
}
