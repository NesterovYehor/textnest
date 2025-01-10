package integrationtests

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/scheduler"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
	testutils "github.com/NesterovYehor/TextNest/services/cleanup_service/tests/unit_tests"
	"github.com/stretchr/testify/assert"
)

func TestExpirationChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute) // Ensure test runs for at most 1 minute
	defer cancel()

	// Initialize configuration
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	defer logFile.Close()

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log := jsonlog.New(multiWriter, slog.LevelInfo)

	// Set up test database
	db, cleanUpDB := testutils.SetupTestDatabase(t, ctx)
	defer cleanUpDB()

	// Set up S3 bucket
	cleanUpS3, err := testutils.SetUpTestS3(ctx)
	assert.NoError(t, err)
	defer cleanUpS3()

	// Kafka setup
	kafkaContainerSetUp, err := container.StartKafka(ctx)
	assert.NoError(t, err)
	defer kafkaContainerSetUp.CleanUp()

	topicName := "example-topic"

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

	// Run expiration processing
	err = srv.ProcessExpirations(ctx)
	assert.NoError(t, err)

	// Test Scheduler
	checker := scheduler.NewChecker(srv, log)
	go func() { // Run scheduler in a separate goroutine to avoid blocking
		checker.Start(ctx, time.Second*10)
	}()

	log.PrintInfo(ctx, "Test completed successfully", nil)
}
