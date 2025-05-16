// Package initializer содержит логику инициализации базы данных.
package initializer

import (
	"context"
	"fmt"
)

type initerRepository interface {
	CreateIfNeededMigrationsTable(ctx context.Context) error
}

type DbInitializerService struct {
	repo initerRepository
}

func New(repo initerRepository) *DbInitializerService {
	return &DbInitializerService{
		repo: repo,
	}
}

// InitDB инициализирует базу данных.
func (s *DbInitializerService) InitDB(ctx context.Context) error {
	err := s.repo.CreateIfNeededMigrationsTable(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize database tables: %w", err)
	}
	return nil
}
