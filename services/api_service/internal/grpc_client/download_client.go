package grpc_clients

import (
	"context"

	paste_download "github.com/NesterovYehor/TextNest/services/api_service/api/download_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DownloadClient struct {
	client paste_download.PasteDownloadClient
	conn   *grpc.ClientConn
}

func NewDownloadClient(target string) (*DownloadClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &DownloadClient{client: paste_download.NewPasteDownloadClient(conn), conn: conn}, nil
}

func (c *DownloadClient) Close() error {
	return c.conn.Close()
}

func (c *DownloadClient) Download(key string) (*paste_download.DownloadResponse, error) {
	req := paste_download.DownloadRequest{
		Key: key,
	}
	resp, err := c.client.Download(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
