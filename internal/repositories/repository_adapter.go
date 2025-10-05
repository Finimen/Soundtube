package repositories

import (
	"context"
	"database/sql"
	"soundtube/pkg"
	"soundtube/pkg/config"
)

type RepositoryAdapter struct {
	db *sql.DB
	*UserRepository
}

func NewRepositoryAdapter(cfg *config.Repository, logger *pkg.CustomLogger) (*RepositoryAdapter, error) {
	var ctx = context.Background()
	var adapter = RepositoryAdapter{}
	var err error
	adapter.db, err = sql.Open("postgress", cfg.Path)
	if err != nil {
		logger.Error("repository initialization completed", err).WithTrace(ctx)
		return nil, err
	}
	if adapter.UserRepository, err = NewUserRepository(adapter.db); err != nil {
		logger.Error("repository initialization completed", err).WithTrace(ctx)
		return nil, err
	}

	logger.Info("repository initialization completed")
	return &adapter, nil
}
