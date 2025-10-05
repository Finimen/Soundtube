package repositories

import (
	"context"
	"soundtube/pkg"
	"time"

	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel/attribute"
)

type TokenBlacklist struct {
	logger *pkg.CustomLogger
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client, logger *pkg.CustomLogger) *TokenBlacklist {
	return &TokenBlacklist{client: client, logger: logger}
}

func (t *TokenBlacklist) Add(ctx context.Context, token string, expiration time.Duration) error {
	ctx, span := t.logger.GetTracer().Start(ctx, "TokenBlacklist.Add")
	defer span.End()

	span.SetAttributes(
		attribute.String("token", token),
	)

	return t.client.Set(formatTokenForList(token), "1", expiration).Err()
}

func (t *TokenBlacklist) Exist(ctx context.Context, token string) (bool, error) {
	exists, err := t.client.Exists(formatTokenForList(token)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func formatTokenForList(token string) string {
	return "bl:" + token
}
