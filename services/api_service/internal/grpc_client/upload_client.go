package grpc_clients

import (
	"context"
	"fmt"

	paste_upload "github.com/NesterovYehor/TextNest/services/api_service/api/upload_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UploadClient struct {
	client paste_upload.PasteUploadClient
	conn   *grpc.ClientConn
}

// Creates a new UploadClient and establishes a connection to the gRPC server
func NewUploadClient(target string) (*UploadClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &UploadClient{client: paste_upload.NewPasteUploadClient(conn), conn: conn}, nil
}

func (c *UploadClient) Close() error {
	return c.conn.Close()
}

// Upload method calls the gRPC Upload RPC
func (c *UploadClient) UploadPaste(ctx context.Context, req *paste_upload.UploadPasteRequest) (string, error) {
	resp, err := c.client.UploadPaste(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.UploadUrl, nil
}

// Upload method calls the gRPC Upload RPC
func (c *UploadClient) ExpireAllUserPastes(ctx context.Context, userId string) (string, error) {
	resp, err := c.client.ExpireAllPastesByUserID(ctx, &paste_upload.ExpireAllPastesByUserIDRequest{
		UserId: userId,
	})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (c *UploadClient) UpdatePaste(
	ctx context.Context,
	key string,
	userId string,
) (string, error) {
	req := &paste_upload.UploadUpdatesRequest{
		Key:    key,
		UserId: userId,
	}

	// Send gRPC request
	resp, err := c.client.UploadUpdates(ctx, req)
	if err != nil {
		return "", fmt.Errorf("gRPC UploadUpdates failed: %w", err)
	}
	return resp.UploadUrl, nil
}
