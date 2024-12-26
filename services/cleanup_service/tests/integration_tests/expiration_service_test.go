package integrationtests

import (
	"context"
	"log"
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
	err := testutils.GetTestEnv()
	assert.NoError(t, err)
	db, cleanUpDB := testutils.SetupTestDatabase(t, ctx)
	defer cleanUpDB()
	cleanUpS3, err := testutils.SetUpTestS3(ctx)
	defer cleanUpS3()
	assert.NoError(t, err)

	opts := &container.KafkaContainerOpts{
		ClusterID:         "test-cluster",
		BrokerPort:        1111,
		Topics:            map[string]int32{"example-topic": 1},
		ReplicationFactor: 1,
	}
	kafkaContainer, brokerAddr, err := container.StartKafka(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to start Kafka: %v", err)
	}

	kafkaCfg := kafka.LoadKafkaConfig([]string{brokerAddr}, []string{"example-tipic"}, "no-group", 1)

	kafkaProd, err := kafka.NewProducer(*kafkaCfg, ctx)
	defer kafkaContainer.Terminate(ctx)

	factory := repository.NewRepositoryFactory(db)
	metadataRepo := factory.CreateMetadataRepository()
	storageRepo, err := factory.CreateStorageRepository()
	srv := services.NewExpirationService(
		metadataRepo, storageRepo,
		kafkaProd,
		testutils.S3TestData.Bucket,
	)

	err = srv.ProcessExpirations(ctx)
	assert.NoError(t, err)
}
