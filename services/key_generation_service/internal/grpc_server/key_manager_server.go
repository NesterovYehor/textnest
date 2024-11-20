package key_manager

import (
	"context"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
	"github.com/redis/go-redis/v9"
)

// KeyManagerServer now includes a redis.Client
type KeyManagerService struct {
	UnimplementedKeyManagerServiceServer
	repo *repository.KeymanagerRepo
}

// NewKeyManagerServer creates a new KeyManagerServer with a Redis client
func NewKeyManagerServer(redisClient *redis.Client, repo *repository.KeymanagerRepo) *KeyManagerService {
	return &KeyManagerService{
		repo: repo,
	}
}

// GetKey now uses the Redis client passed in the server struct
func (s *KeyManagerService) GetKey(ctx context.Context, req *GetKeyRequest) (*GetKeyResponse, error) {
	key, err := s.repo.GetKey() // Pass Redis client to GetKey
	if err != nil {
		return &GetKeyResponse{Error: err.Error()}, nil
	}
	return &GetKeyResponse{Key: key}, nil
}
