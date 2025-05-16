package checker

import (
	"context"
	"fmt"
	"time"

	"migrator/internal/entity"
)

type migratorSrv interface {
	ApplyMigration(ctx context.Context, migrationIDs []int64, userID int64) (time.Time, error)
	CreateMigration(ctx context.Context, name string, description string, script string, rollbackScript string, userID int64) (int64, error)
	GetMigration(ctx context.Context, migrationID int64) (entity.MigrationInfo, error)
	ListMigrations(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error)
	RollbackMigration(ctx context.Context, migrationID int64, userID int64) (time.Time, error)
}

type authClient interface {
	CheckPermissionApply(ctx context.Context, userID int64) (bool, error)
	CheckPermissionApplyOther(ctx context.Context, userID int64) (bool, error)
	CheckPermissionCreate(ctx context.Context, userID int64) (bool, error)
	CheckPermissionGet(ctx context.Context, userID int64) (bool, error)
	CheckPermissionList(ctx context.Context, userID int64) (bool, error)
	CheckPermissionRollback(ctx context.Context, userID int64) (bool, error)
	CheckPermissionRollbackOther(ctx context.Context, userID int64) (bool, error)
}

// MigratorWithAuth is a wrapper around Migrator that adds authorization checks.
type MigratorWithAuth struct {
	migrator   migratorSrv
	authClient authClient
}

// NewMigratorWithAuth creates a new MigratorWithAuth.
func NewMigratorWithAuth(migrator migratorSrv, authClient authClient) *MigratorWithAuth {
	return &MigratorWithAuth{
		migrator:   migrator,
		authClient: authClient,
	}
}

// CreateMigration creates a new migration after checking permissions.
func (mwa *MigratorWithAuth) CreateMigration(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error) {
	hasPermission, err := mwa.authClient.CheckPermissionCreate(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("auth check failed for CreateMigration: %w", err)
	}
	if !hasPermission {
		return 0, fmt.Errorf("%w: user %d lacks permission to create migration", entity.ErrPermissionDenied, userID)
	}

	return mwa.migrator.CreateMigration(ctx, name, description, script, rollbackScript, userID)
}

// ApplyMigration applies migrations after checking permissions.
func (mwa *MigratorWithAuth) ApplyMigration(ctx context.Context, migrationIDs []int64, userID int64) (time.Time, error) {
	hasPermission, err := mwa.authClient.CheckPermissionApply(ctx, userID)
	if err != nil {
		return time.Time{}, fmt.Errorf("auth check failed for ApplyMigration: %w", err)
	}
	if !hasPermission {
		return time.Time{}, fmt.Errorf("%w: user %d lacks permission to apply migrations", entity.ErrPermissionDenied, userID)
	}

	return mwa.migrator.ApplyMigration(ctx, migrationIDs, userID)
}

// RollbackMigration откатывает миграцию.
func (mwa *MigratorWithAuth) RollbackMigration(ctx context.Context, migrationID, actorUserID int64) (time.Time, error) {
	migrationInfo, err := mwa.GetMigration(ctx, migrationID)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get migration info for rollback auth check: %w", err)
	}

	creatorUserID := migrationInfo.CreatedBy

	var hasPermission bool
	var permCheckErr error
	var permType string

	if creatorUserID == actorUserID {
		permType = "PERMISSION_ROLLBACK (own)"
		hasPermission, permCheckErr = mwa.authClient.CheckPermissionRollback(ctx, actorUserID)
	} else {
		permType = "PERMISSION_ROLLBACK_OTHER"
		hasPermission, permCheckErr = mwa.authClient.CheckPermissionRollbackOther(ctx, actorUserID)
	}

	if permCheckErr != nil {
		return time.Time{}, fmt.Errorf("auth check failed for RollbackMigration (%s): %w", permType, permCheckErr)
	}
	if !hasPermission {
		return time.Time{}, fmt.Errorf("%w: user %d lacks %s for migration %d", entity.ErrPermissionDenied, actorUserID, permType, migrationID)
	}

	return mwa.migrator.RollbackMigration(ctx, migrationID, actorUserID)
}

// ListMigrations возвращает список миграций.
func (mwa *MigratorWithAuth) ListMigrations(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error) {
	return mwa.migrator.ListMigrations(ctx, statusFilter)
}

// GetMigration возвращает миграцию по ее ID.
func (mwa *MigratorWithAuth) GetMigration(ctx context.Context, migrationID int64) (entity.MigrationInfo, error) {
	return mwa.migrator.GetMigration(ctx, migrationID)
}
