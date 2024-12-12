package integrationtests

import (
	"context"
	"testing"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
	testutils "github.com/NesterovYehor/TextNest/services/key_generation_service/tests/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	client, cleanup := testutils.StartTestRedis(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(client)
	repo.FillKeys(2)

	key, err := repo.GetKey(ctx)
	require.NoError(t, err, "Failed to get key")
	require.NotEmpty(t, key, "Key is empty")

	err = repo.ReallocateKey(key)
	assert.NoError(t, err, "Failed to reallocate key")
}
