package intiter

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

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

type Repository struct {
	conn Excecutor
}

func New(conn Excecutor) *Repository {
	return &Repository{
		conn: conn,
	}
}

const createMigrationsTableQuery = `
CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    script TEXT NOT NULL,
    rollback_script TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status_updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
`

// CreateIfNeededMigrationsTable создает таблицу миграций, если ее нет.
func (r *Repository) CreateIfNeededMigrationsTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createMigrationsTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}
