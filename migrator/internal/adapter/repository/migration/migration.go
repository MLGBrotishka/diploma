package migration

import (
	"context"
	"fmt"
	"time"

	"migrator/internal/entity"
	"migrator/pkg/logger"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

type migrationRepository interface {
	Get(ctx context.Context, migrationID int64) (entity.MigrationInfo, error)
	Create(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error)
	Apply(ctx context.Context, script string) error
	SetStatus(ctx context.Context, migrationID int64, updatedAt time.Time, status entity.MigrationStatus) error
	List(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error)
	GetLatestAppliedMigration(ctx context.Context) (entity.MigrationInfo, error)
	DoInTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

var _ migrationRepository = (*Repository)(nil)

// Excecutor - интерфейс для выполнения запросов на базе данных.
type Excecutor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type (
	txKey struct{}
)

type Repository struct {
	conn Excecutor
}

func New(conn Excecutor) *Repository {
	return &Repository{
		conn: conn,
	}
}

const getQuery = `-- Get
	SELECT 
		id, 
		name, 
		description, 
		script, 
		rollback_script, 
		status, 
		created_by, 
		status_updated_at
	FROM migrations
	WHERE id = $1
`

func (r *Repository) Get(ctx context.Context, migrationID int64) (entity.MigrationInfo, error) {
	var migration entity.MigrationInfo
	err := r.Do(ctx).QueryRow(ctx, getQuery, migrationID).
		Scan(
			&migration.ID,
			&migration.Name,
			&migration.Description,
			&migration.Script,
			&migration.RollbackScript,
			&migration.Status,
			&migration.CreatedBy,
			&migration.StatusUpdatedAt,
		)
	if err != nil {
		return entity.MigrationInfo{}, fmt.Errorf("get migration: %w", err)
	}
	return migration, nil
}

const createQuery = `-- Create
	INSERT INTO migrations (name, description, script, rollback_script, created_by, status, created_at, status_updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
	RETURNING id
`

func (r *Repository) Create(ctx context.Context, name, description, script, rollbackScript string, userID int64) (int64, error) {
	var id int64
	err := r.Do(ctx).QueryRow(
		ctx,
		createQuery,
		name,
		description,
		script,
		rollbackScript,
		userID,
		entity.StatusPending,
		time.Now().UTC(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create migration: %w", err)
	}

	return id, nil
}

func (r *Repository) Apply(ctx context.Context, script string) error {
	_, err := r.Do(ctx).Exec(ctx, script)
	if err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}
	return nil
}

const setStatusQuery = `-- SetStatus
	UPDATE migrations
	SET status = $1, status_updated_at = $2
	WHERE id = $3
`

func (r *Repository) SetStatus(ctx context.Context, migrationID int64, updatedAt time.Time, status entity.MigrationStatus) error {
	_, err := r.Do(ctx).Exec(ctx, setStatusQuery, status, updatedAt, migrationID)
	if err != nil {
		return fmt.Errorf("set status: %w", err)
	}
	return nil
}

const listQuery = `-- List
	SELECT
		id,
		name,
		description,
		script,
		rollback_script,
		status,
		created_by,
		status_updated_at
	FROM migrations
	WHERE status = $1 OR $1 = ''
	ORDER BY id
`

func (r *Repository) List(ctx context.Context, statusFilter string) ([]entity.MigrationInfo, error) {
	rows, err := r.Do(ctx).Query(ctx, listQuery, statusFilter)
	if err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}
	defer rows.Close()

	var migrations []entity.MigrationInfo
	for rows.Next() {
		var migration entity.MigrationInfo
		err := rows.Scan(
			&migration.ID,
			&migration.Name,
			&migration.Description,
			&migration.Script,
			&migration.RollbackScript,
			&migration.Status,
			&migration.CreatedBy,
			&migration.StatusUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan migration: %w", err)
		}
		migrations = append(migrations, migration)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return migrations, nil
}

const getLatestAppliedMigrationQuery = `-- GetLatestAppliedMigration
	SELECT
		id,
		name,
		description,
		script,
		rollback_script,
		status,
		created_by,
		status_updated_at
	FROM migrations
	WHERE status = $1
	ORDER BY id DESC
	LIMIT 1
`

func (r *Repository) GetLatestAppliedMigration(ctx context.Context) (entity.MigrationInfo, error) {
	var migration entity.MigrationInfo
	err := r.Do(ctx).QueryRow(ctx, getLatestAppliedMigrationQuery, entity.StatusApplied).
		Scan(
			&migration.ID,
			&migration.Name,
			&migration.Description,
			&migration.Script,
			&migration.RollbackScript,
			&migration.Status,
			&migration.CreatedBy,
			&migration.StatusUpdatedAt,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.MigrationInfo{}, nil
		}
		return entity.MigrationInfo{}, fmt.Errorf("get latest applied migration: %w", err)
	}

	return migration, nil
}

func (r *Repository) DoInTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := r.Do(ctx).Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.Error(fmt.Errorf("transaction rollback error: %w", err))
		}
	}()

	ctx = context.WithValue(ctx, txKey{}, tx)

	err = f(ctx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) Do(ctx context.Context) Excecutor {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)

	if ok && isValidTx(tx) {
		return tx
	}

	return r.conn
}

func isValidTx(tx pgx.Tx) bool {
	if tx == nil {
		return false
	}

	if tx.Conn() == nil || tx.Conn().IsClosed() {
		return false
	}

	return true
}
