// Package entity определяет сущности, используемые в сервисе миграций.
package entity

import "time"

type MigrationInfo struct {
	ID              int64           `json:"id" db:"id"`
	Name            string          `json:"name" db:"name"`
	Description     string          `json:"description" db:"description"`
	Script          string          `json:"script" db:"script"`
	RollbackScript  string          `json:"rollback_script" db:"rollback_script"`
	Status          MigrationStatus `json:"status" db:"status"`
	CreatedBy       int64           `json:"created_by" db:"created_by"`
	StatusUpdatedAt time.Time       `json:"status_updated_at" db:"status_updated_at"`
}

type MigrationStatus string

const (
	StatusPending    MigrationStatus = "pending"
	StatusApplied    MigrationStatus = "applied"
	StatusRolledBack MigrationStatus = "rolled_back"
)

func (s MigrationStatus) String() string {
	return string(s)
}
