package container

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainerOpts struct {
	Addr int
}

type RedisContainer struct {
	*redis.RedisContainer
}

func StartRedis(ctx context.Context, opts *RedisContainerOpts) (*RedisContainer, error) {
	redisContainer, err := redis.Run(ctx,
		"redis:7",
		testcontainers.WithHostPortAccess(opts.Addr),
	)
	if err != nil {
		return nil, err
	}
	return &RedisContainer{redisContainer}, nil
}
