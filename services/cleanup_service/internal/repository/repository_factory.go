package repository

import (
	"database/sql"
)

// RepositoryFactory is responsible for creating repositories.
type RepositoryFactory struct {
	db *sql.DB
}

// NewRepositoryFactory creates a new instance of the factory with dependencies.
func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// CreateMetadataRepository initializes a MetadataRepository.
func (f *RepositoryFactory) CreateMetadataRepository() MetadataRepository {
	return newDBRepository(f.db)
}

// CreateStorageRepository initializes a StorageRepository.
func (f *RepositoryFactory) CreateStorageRepository() (StorageRepository, error) {
	return newS3Storage()
}
