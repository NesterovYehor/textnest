package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
	"github.com/stretchr/testify/assert"
)

func DownloaServieTest(t *testing.T) {
	// Create a context for container setup and tests
	ctx := context.Background()

	// Start Postgres container
	postgresContainer, err := container.StartPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to start Postgres container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)

	// Start Redis container
	redisContainer, err := container.StartRedis(ctx, &container.RedisContainerOpts{Addr: 6379})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}
	defer redisContainer.Terminate(ctx)

	// Start Kafka container
	kafkaContainer, err := container.StartKafka(ctx)
	if err != nil {
		t.Fatalf("Failed to start Kafka container: %v", err)
	}
	defer kafkaContainer.Terminate(ctx)

	// Setup the database connection
	db, err := postgresContainer.GetConnection(ctx)
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}

	// Setup Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Setup Minio S3 (or real S3 setup)
	s3Client := setupRealS3()

	// Insert test metadata directly into the database
	_, err = db.ExecContext(ctx, `
		INSERT INTO metadata (key, created_at, expiration_date) 
		VALUES ($1, $2, $3)
	`, "testkey", time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		t.Fatalf("Failed to insert test metadata: %v", err)
	}

	// Insert test content into S3 (real storage)
	err = insertTestContentToS3(s3Client, "testkey", "test content data")
	if err != nil {
		t.Fatalf("Failed to insert test content to S3: %v", err)
	}

	// Setup the Cache and Kafka producer
	metadataCache := cache.NewRedisCache(ctx, "localhost:6379")
	contentCache := cache.NewRedisCache(ctx, "localhost:6379")
	kafkaProducer, err := kafka.NewProducer(kafka.KafkaConfig{}, ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	// Create the DownloadService
	metadataRepo := repository.NewMetadataRepo(db)
	storageRepo := repository.NewStorageRepository() // Assuming you have a real storage repo

	downloadService, err := services.NewDownloadService(
		storageRepo,
		metadataRepo,
		log,
		ctx,
		"localhost:6379",
		"localhost:6379",
		kafka.KafkaConfig{},
	)
	if err != nil {
		t.Fatalf("Failed to create DownloadService: %v", err)
	}

	// Test DownloadService Download method
	t.Run("Download valid paste", func(t *testing.T) {
		// Directly invoke the Download method (no gRPC)
		result, err := downloadService.Download(ctx, "testkey")
		assert.NoError(t, err)
		assert.Equal(t, "testkey", result.Key)
		assert.Equal(t, "test content data", result.Content) // Ensure content matches
		assert.NotNil(t, result.ExpirationDate)
	})

	t.Run("Download invalid paste", func(t *testing.T) {
		// Directly invoke the Download method (no gRPC)
		result, err := downloadService.Download(ctx, "invalidkey")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
