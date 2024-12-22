package integrationtests

import (
	"context"
	"testing"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/services"
	pb "github.com/NesterovYehor/TextNest/services/key_generation_service/proto"
	testutils "github.com/NesterovYehor/TextNest/services/key_generation_service/tests/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestKeyGenerationService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Start Redis and set up cleanup
	client, cleanup := testutils.StartTestRedis(t, ctx)
	defer cleanup()

	// Initialize repository and service
	repo := repository.NewRepository(client)
	assert.NotNil(t, repo, "Failed to initialize repository")

	// Pre-fill keys
	repo.FillKeys(1) // Assume FillKeys adds at least one key

	service := services.NewKeyManagerServer(repo)
	assert.NotNil(t, service, "Failed to initialize KeyManagerService")

	// Test GetKey
	res, err := service.GetKey(ctx, &pb.GetKeyRequest{})
	assert.NoError(t, err, "GetKey returned an error")
	assert.NotNil(t, res, "Response should not be nil")
	assert.NotEmpty(t, res.Key, "Key in the response should not be empty")
}
