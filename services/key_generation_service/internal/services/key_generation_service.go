package services

import (
	"context"

	pb "github.com/NesterovYehor/TextNest/services/key_generation_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
)

// KeyManagerServer now includes a redis.Client
type KeyManagerService struct {
	pb.UnimplementedKeyManagerServiceServer
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
