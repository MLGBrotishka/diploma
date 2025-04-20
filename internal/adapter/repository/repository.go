package repository

import (
	"context"
	"fmt"

	"migrator/pkg/logger"

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

type (
	txKey  struct{}
	txFunc func(ctx context.Context) error
)

type Repository struct {
	conn Excecutor
}

func New(conn Excecutor) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) DoInTransaction(ctx context.Context, f txFunc) error {
	tx, err := r.do(ctx).Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.Error(fmt.Errorf("transaction rollback error: %w", err))
		}
	}()

	err = f(ctx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) do(ctx context.Context) Excecutor {
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
