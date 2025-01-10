package integrationtests

import (
	"context"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
	testutils "github.com/NesterovYehor/TextNest/services/cleanup_service/tests/unit_tests"
	"github.com/stretchr/testify/assert"
)

func TestProcessExpirations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up test database
	db, cleanUpDB := testutils.SetupTestDatabase(t, ctx)
	defer cleanUpDB()

	// Set up S3 bucket
	cleanUpS3, err := testutils.SetUpTestS3(ctx)
	assert.NoError(t, err)
	defer cleanUpS3()

	// Kafka options
	topicName := "example-topic"

	// Start Kafka container
	kafkaContainerSetUp, err := container.StartKafka(ctx)
	assert.NoError(t, err)
	defer kafkaContainerSetUp.CleanUp()

	// Configure Kafka producer
	kafkaCfg := kafka.LoadKafkaConfig(kafkaContainerSetUp.BrokerAddr, []string{topicName}, "no-group", 1)
	kafkaProd, err := kafka.NewProducer(*kafkaCfg, ctx)
	assert.NoError(t, err)

	// Create repositories and services
	factory := repository.NewRepositoryFactory(db)
	metadataRepo := factory.CreateMetadataRepository()
	storageRepo, err := factory.CreateStorageRepository()
	assert.NoError(t, err)

	// Expiration service
	srv := services.NewExpirationService(
		metadataRepo, storageRepo,
		kafkaProd,
		testutils.S3TestData.Bucket,
	)

	// Execute expiration processing
	err = srv.ProcessExpirations(ctx)
	assert.NoError(t, err)
}
