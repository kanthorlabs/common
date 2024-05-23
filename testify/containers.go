package testify

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func RedisContainer(ctx context.Context) (*redis.RedisContainer, error) {
	return redis.RunContainer(
		ctx,
		testcontainers.WithImage("redis:7-alpine"),
	)
}
