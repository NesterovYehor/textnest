package services

import (
	"context"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
	pb "github.com/NesterovYehor/TextNest/services/key_generation_service/proto"
)

// KeyManagerServer now includes a redis.Client
type KeyManagerService struct {
	pb.UnimplementedKeyGeneratorServer
	repo *repository.KeyGeneratorRepository
}

// NewKeyManagerServer creates a new KeyManagerServer with a Redis client
func NewKeyManagerServer(repo *repository.KeyGeneratorRepository) *KeyManagerService {
	return &KeyManagerService{
		repo: repo,
	}
}

// GetKey now uses the Redis client passed in the server struct
func (s *KeyManagerService) GetKey(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	key, err := s.repo.GetKey(ctx) // Pass Redis client to GetKey
	if err != nil {
		return &pb.GetKeyResponse{Error: err.Error()}, nil
	}
	return &pb.GetKeyResponse{Key: key}, nil
}
