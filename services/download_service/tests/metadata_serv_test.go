package tests

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestFetchMetadataByKey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	connString, cleanupRedis := SetUpRedis(ctx, t)
	defer cleanupRedis()

	db, cleanupPg := SetUpPostgres(ctx, t)
	defer cleanupPg()
	query := `
        INSERT INTO metadata(key, title, user_id, expiration_date) 
        VALUES ($1, NULLIF($2, ''), $3, $4)
    `
	_, err := db.ExecContext(ctx, query, key, title, userId, expirationDate.AsTime())
	if err != nil {
		t.Fatal("Failed insert test data to postgres test container")
	}

	cache, err := cache.NewRedisCache(connString)
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewMetadataRepo(db)
	kafkaProd := SetUpKafka(ctx, t)

	srv := services.NewFetchMetadataService(repo, cache, kafkaProd)
	res, err := srv.FetchMetadataByKey(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, title, res.Title)


}
