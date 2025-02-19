package tests

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestDownloadMetadata(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	db, cleanup := SetUpPostgres(ctx, t)
	defer cleanup()

	query := `
        INSERT INTO metadata(key, title, user_id, expiration_date) 
        VALUES ($1, NULLIF($2, ''), $3, $4)
    `
	_, err := db.ExecContext(ctx, query, key, title, userId, expirationDate.AsTime())
	assert.NoError(t, err)

	repo := repository.NewMetadataRepo(db)
	res, err := repo.DownloadPasteMetadata(ctx, key)
	assert.NoError(t, err)

	assert.Equal(t, key, res.Key)
	assert.WithinDuration(t, expirationDate.AsTime(), res.ExpiredDate.AsTime(), time.Second)
	assert.Equal(t, title, res.Title)
}
