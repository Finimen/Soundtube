package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"soundtube/pkg"
	"soundtube/pkg/config"
	"time"

	_ "github.com/lib/pq"
)

type RepositoryAdapter struct {
	db *sql.DB
	*UserRepository
	*SoundRepository
}

func NewRepositoryAdapter(dbCfg *config.Database, connCfg *config.DatabaseConnections, logger *pkg.CustomLogger) (*RepositoryAdapter, error) {
	var ctx = context.Background()
	var adapter = RepositoryAdapter{}
	var err error
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)
	adapter.db, err = sql.Open("postgres", conn)
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

	if adapter.UserRepository, err = NewUserRepository(adapter.db, logger); err != nil {
		logger.Error("user repository failed", err).WithTrace(ctx)
		return nil, err
	}

	if adapter.SoundRepository, err = NewSoundRepository(adapter.db, logger); err != nil {
		logger.Error("sound repository failed", err).WithTrace(ctx)
		return nil, err
	}

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
