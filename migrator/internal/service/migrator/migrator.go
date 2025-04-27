package migrator

import (
	"context"
	"fmt"
	"time"

	"migrator/internal/entity"
)

type MigrationService interface {
	CreateMigration(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error)
	ApplyMigration(ctx context.Context, migrationIDs []int64, userID int64) (time.Time, error)
	RollbackMigration(ctx context.Context, migrationID, userID int64) (time.Time, error)
	ListMigrations(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error)
	GetMigration(ctx context.Context, migrationID int64) (entity.MigrationInfo, error)
}

var _ MigrationService = (*Migrator)(nil)

type migrationRepository interface {
	Get(ctx context.Context, migrationID int64) (entity.MigrationInfo, error)
	Create(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error)
	Apply(ctx context.Context, script string) error
	SetStatus(ctx context.Context, migrationID int64, updatedAt time.Time, status entity.MigrationStatus) error
	List(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error)
	GetLatestAppliedMigration(ctx context.Context) (entity.MigrationInfo, error)
	DoInTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

type Migrator struct {
	repo migrationRepository
}

func New(repo migrationRepository) *Migrator {
	return &Migrator{
		repo: repo,
	}
}

// CreateMigration создает новую миграцию.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	name: string - Название миграции.
//	description: string - Описание миграции.
//	script: string - Текст скрипта миграции.
//	rollbackScript: string - Текст скрипта отката миграции.
//	userID: int64 - Идентификатор пользователя, создающего миграцию.
//
// Возвращает:
//
//	int64: Уникальный идентификатор созданной миграции.
//	error: Ошибка, если таковая имеется.
func (m *Migrator) CreateMigration(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error) {
	migrationID, err := m.repo.Create(ctx, name, description, script, rollbackScript, userID)
	if err != nil {
		return 0, err
	}

	return migrationID, nil
}

// ApplyMigration применяет миграции.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	migrationIDs: []int64 - Уникальные идентификаторы миграций в соответствии с порядком применения.
//	userID: int64 - Идентификатор пользователя, применяющего миграцию.
//
// Возвращает:
//
//	time.Time: Дата и время применения миграции.
//	error: Ошибка, если таковая имеется.
func (m *Migrator) ApplyMigration(ctx context.Context, migrationIDs []int64, userID int64) (time.Time, error) {
	var appliedAt time.Time

	err := m.repo.DoInTransaction(ctx, func(ctx context.Context) error {
		var migrations []entity.MigrationInfo
		for _, migrationID := range migrationIDs {
			migration, err := m.repo.Get(ctx, migrationID)
			if err != nil {
				return fmt.Errorf("m.repo.GetMigration: %w", err)
			}

			if migration.Status != entity.StatusPending {
				return fmt.Errorf("migration %d is not pending", migrationID)
			}

			migrations = append(migrations, migration)
		}

		for _, migration := range migrations {
			err := m.repo.Apply(ctx, migration.Script)
			if err != nil {
				return fmt.Errorf("m.repo.ApplyMigration: %w", err)
			}

			err = m.repo.SetStatus(ctx, migration.ID, time.Now(), entity.StatusApplied)
			if err != nil {
				return fmt.Errorf("m.repo.SetApplied: %w", err)
			}
		}

		appliedAt = time.Now()
		return nil
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("m.repo.DoInTransaction: %w", err)
	}

	return appliedAt, nil
}

// RollbackMigration откатывает миграцию.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	migrationID: int64 - Уникальный идентификатор миграции для отката.
//	userID: int64 - Идентификатор пользователя, выполняющего откат.
//
// Возвращает:
//
//	time.Time: Дата и время отката миграции.
//	error: Ошибка, если таковая имеется.
func (m *Migrator) RollbackMigration(ctx context.Context, migrationID, userID int64) (time.Time, error) {
	var rolledBackAt time.Time

	err := m.repo.DoInTransaction(ctx, func(ctx context.Context) error {
		migration, err := m.repo.Get(ctx, migrationID)
		if err != nil {
			return fmt.Errorf("m.repo.GetMigration: %w", err)
		}

		if migration.Status != entity.StatusApplied {
			return fmt.Errorf("migration %d is not pending", migrationID)
		}

		latestAppliedMigration, err := m.repo.GetLatestAppliedMigration(ctx)
		if err != nil {
			return fmt.Errorf("m.repo.GetLatestAppliedMigration: %w", err)
		}

		if latestAppliedMigration.ID != migrationID {
			return fmt.Errorf("not last migration")
		}

		err = m.repo.Apply(ctx, migration.RollbackScript)
		if err != nil {
			return fmt.Errorf("m.db.ApplyMigration: %w", err)
		}

		rolledBackAt = time.Now()

		err = m.repo.SetStatus(ctx, migration.ID, rolledBackAt, entity.StatusRolledBack)
		if err != nil {
			return fmt.Errorf("m.repo.SetStatus: %w", err)
		}

		return nil
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("m.repo.DoInTransaction: %w", err)
	}

	return rolledBackAt, nil
}

// ListMigrations возвращает список миграций.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	statusFilter: string - Фильтр по статусу миграции.
//
// Возвращает:
//
//	[]entity.MigrationInfo: Список информации о миграциях.
//	error: Ошибка, если таковая имеется.
func (m *Migrator) ListMigrations(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error) {
	migrations, err := m.repo.List(ctx, statusFilter)
	if err != nil {
		return nil, fmt.Errorf("m.repo.ListMigrations: %w", err)
	}

	return migrations, nil
}

// GetMigration возвращает миграцию по ее ID.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	migrationID: int64 - Уникальный идентификатор миграции.
//
// Возвращает:
//
//	entity.MigrationInfo: Информация о миграции.
//	error: Ошибка, если таковая имеется.
func (m *Migrator) GetMigration(ctx context.Context, migrationID int64) (entity.MigrationInfo, error) {
	migration, err := m.repo.Get(ctx, migrationID)
	if err != nil {
		return entity.MigrationInfo{}, fmt.Errorf("m.repo.GetMigration: %w", err)
	}

	return migration, nil
}
