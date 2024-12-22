package services

import (
	"context"
	"database/sql"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
)

type ServiceFactory struct {
	cfg *config.Config
	log *jsonlog.Logger
}

func NewServiceFactory(cfg *config.Config, log *jsonlog.Logger) *ServiceFactory {
	return &ServiceFactory{cfg: cfg, log: log}
}

func (f *ServiceFactory) CreateFetchMetadataService(ctx context.Context) (*FetchMetadataService, error) {
	db, err := f.openDB()
	if err != nil {
		return nil, err
	}
	metadataRepo := repository.NewMetadataRepo(db)
	kafkaProducer, err := kafka.NewProducer(f.cfg.Kafka, ctx)
	if err != nil {
		return nil, err
	}
	cache := cache.NewRedisCache(ctx, f.cfg.RedisMetadataAddr)

	return newFetchMetadataService(metadataRepo, cache, f.log, kafkaProducer)
}

func (f *ServiceFactory) CreateFetchContentService(ctx context.Context) (*FetchContentService, error) {
	storageRepo, err := repository.NewStorageRepository(f.cfg.BucketName)
	if err != nil {
		return nil, err
	}
	cache := cache.NewRedisCache(ctx, f.cfg.RedisContentAddr)

	return newFetchContentService(storageRepo, cache, f.log, f.cfg.BucketName)
}

func (f *ServiceFactory) openDB() (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("postgres", f.cfg.DBURL)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
