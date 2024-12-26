package grpc_clients

import (
	"context"
	"time"

	paste_upload "github.com/NesterovYehor/TextNest/services/api_service/api/upload_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UploadClient struct {
	client paste_upload.PasteUploadClient
}

// Creates a new UploadClient and establishes a connection to the gRPC server
func NewUploadClient(target string) (*UploadClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &UploadClient{client: paste_upload.NewPasteUploadClient(conn)}, nil
}

// Upload method calls the gRPC Upload RPC
func (c *UploadClient) Upload(key string, expirationDate time.Time, data []byte) (string, error) {
	resp, err := c.client.Upload(context.Background(), &paste_upload.UploadRequest{
		Key:            key,
		ExpirationDate: timestamppb.New(expirationDate),
		Data:           data,
	})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
