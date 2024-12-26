package grpc_clients

import (
	"context"
	"sync"

	key_generation "github.com/NesterovYehor/TextNest/services/api_service/api/key_generation_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	clientOnce      sync.Once
	keyGenClient    *KeyGeneratorClient
	clientInitError error
)

// KeyGeneratorClient wraps the generated gRPC client for easier use.
type KeyGeneratorClient struct {
	client key_generation.KeyGeneratorClient
}

// NewKeyGeneratorClient creates a new KeyGeneratorClient and establishes a connection to the given target.
func NewKeyGeneratorClient(target string) (*KeyGeneratorClient, error) {
	clientOnce.Do(func() {
		// Use grpc.Dial to create a connection to the server.
		conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			clientInitError = err
			return
		}
		// Return a new KeyGeneratorClient initialized with the gRPC client.
		keyGenClient = &KeyGeneratorClient{
			client: key_generation.NewKeyGeneratorClient(conn),
		}
	})
	return keyGenClient, clientInitError
}

// GetKey sends a request to the Key Generator service to fetch a key.
func (c *KeyGeneratorClient) GetKey(ctx context.Context) (string, error) {
	// Perform the gRPC call with the provided context.
	resp, err := c.client.GetKey(ctx, &key_generation.GetKeyRequest{})
	if err != nil {
		return "", err // Return error if the call fails.
	}

	// Return the key from the response.
	return resp.Key, nil
}
