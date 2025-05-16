package auth

import (
	"context"
	"fmt"

	"auth/internal/entity"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

// GetUserByLogin retrieves a user by their login.
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (entity.User, error) {
	query := `SELECT id, login, password_hash, created_at, updated_at, is_active FROM users WHERE login = $1 AND is_active = TRUE`
	var user entity.User
	err := r.conn.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PassHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err == pgx.ErrNoRows {
		return entity.User{}, fmt.Errorf("user not found: %w", err) // Wrap ErrNoRows
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user by login: %w", err)
	}
	return user, nil
}

// SaveUser saves a new user to the database.
func (r *Repository) SaveUser(ctx context.Context, login string, passwordHash []byte) (int64, error) {
	query := `INSERT INTO users (login, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	var userID int64
	err := r.conn.QueryRow(ctx, query, login, passwordHash).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	return userID, nil
}

// CheckUserPermission checks if a user has a specific permission.
func (r *Repository) CheckUserPermission(ctx context.Context, userID int64, permission entity.Permission) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM users u
            JOIN user_roles ur ON u.id = ur.user_id
            JOIN roles r ON ur.role_id = r.id
            JOIN role_permissions rp ON r.id = rp.role_id
            JOIN permissions p ON rp.permission_id = p.id
            WHERE u.id = $1 AND p.name = $2 AND u.is_active = TRUE
        )
    `
	var hasPermission bool
	err := r.conn.QueryRow(ctx, query, userID, permission.String()).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}
	return hasPermission, nil
}

// SetUserInactive sets a user's is_active status to false.
func (r *Repository) SetUserInactive(ctx context.Context, userID int64) error {
	query := `UPDATE users SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to set user inactive: %w", err)
	}
	return nil
}
