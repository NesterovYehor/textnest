package tests

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheAPI(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	redisAddr, cleanup := SetUpRedis(ctx, t)
	defer cleanup()
	cacheClient, err := cache.NewRedisCache(redisAddr)
	assert.NoError(t, err)
	assert.NoError(t, cacheClient.Set(ctx, key, &data))

	res, found, err := cacheClient.Get(ctx, key)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, key, res.Key)
	assert.Equal(t, title, res.Title)

	assert.NoError(t, cacheClient.Delete(ctx, key))
	assert.NoError(t, cacheClient.Clear(ctx))
	assert.NoError(t, cacheClient.Close())
}
