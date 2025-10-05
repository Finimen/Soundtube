package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"soundtube/pkg"
	"soundtube/pkg/config"
	"time"

	"github.com/go-redis/redis"
)

type RepositoryAdapter struct {
	db *sql.DB
	*UserRepository
	*TokenBlacklist
}

func NewRepositoryAdapter(dbCfg *config.Database, connCfg *config.DatabaseConnections, client *redis.Client, logger *pkg.CustomLogger) (*RepositoryAdapter, error) {
	var ctx = context.Background()
	var adapter = RepositoryAdapter{}
	var err error
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)
	adapter.db, err = sql.Open("postgress", conn)
	if err != nil {
		logger.Error("repository initialization completed", err).WithTrace(ctx)
		return nil, err
	}

	adapter.db.SetMaxOpenConns(connCfg.MaxOpenConns)
	adapter.db.SetMaxIdleConns(connCfg.MaxIdleConns)
	adapter.db.SetConnMaxIdleTime(time.Duration(connCfg.ConnMaxIdleTime) * time.Minute)
	adapter.db.SetConnMaxLifetime(time.Duration(connCfg.ConnMaxLifetime) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = adapter.db.PingContext(ctx); err != nil {
		return nil, err
	}

	if adapter.UserRepository, err = NewUserRepository(adapter.db); err != nil {
		logger.Error("repository initialization completed", err).WithTrace(ctx)
		return nil, err
	}

	adapter.TokenBlacklist = NewTokenBlacklist(client, logger)

	logger.Info("repository initialization completed")
	return &adapter, nil
}

func (r *RepositoryAdapter) HealthCheck(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryAdapter) Close() error {
	if err := r.db.Close(); err != nil {
		return err
	}

	return nil
}
