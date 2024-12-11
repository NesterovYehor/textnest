package container

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainerOpts struct {
	addr int
}

type RedisContainer struct {
	*redis.RedisContainer
}

func StartRedis(ctx context.Context, opts *RedisContainerOpts, t *testing.T) (*RedisContainer, error) {
	redisContainer, err := redis.Run(ctx,
		"redis:7",
		testcontainers.WithHostPortAccess(opts.addr),
	)
	if err != nil {
		return nil, err
	}
	return &RedisContainer{redisContainer}, nil
}
