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

const createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE
);
`

// CreateIfNeededUsersTable создает таблицу пользователей, если ее нет.
func (r *Repository) CreateIfNeededUsersTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createUsersTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

const createRolesTableQuery = `
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);
`

// CreateIfNeededRolesTable creates the roles table if it doesn't exist.
func (r *Repository) CreateIfNeededRolesTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createRolesTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}
	return nil
}

const createPermissionsTableQuery = `
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);
`

// CreateIfNeededPermissionsTable creates the permissions table if it doesn't exist.
func (r *Repository) CreateIfNeededPermissionsTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createPermissionsTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create permissions table: %w", err)
	}
	return nil
}

const createRolePermissionsTableQuery = `
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);
`

// CreateIfNeededRolePermissionsTable creates the role_permissions table if it doesn't exist.
func (r *Repository) CreateIfNeededRolePermissionsTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createRolePermissionsTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create role_permissions table: %w", err)
	}
	return nil
}

const createUserRolesTableQuery = `
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);
`

// CreateIfNeededUserRolesTable creates the user_roles table if it doesn't exist.
func (r *Repository) CreateIfNeededUserRolesTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createUserRolesTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create user_roles table: %w", err)
	}
	return nil
}
