package grpc_clients

import (
	"context"
	"time"

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

func (c *DownloadClient) DownloadByKey(key string) (*paste_download.DownloadByKeyResponse, error) {
	req := paste_download.DownloadByKeyRequest{
		Key: key,
	}
	resp, err := c.client.DownloadByKey(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *DownloadClient) DownloadByUserId(userId string, limit, offset int32) ([]*paste_download.Metadata, error) {
	req := paste_download.DownloadByUserIdRequest{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.DownloadByUserId(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.Objects, nil
}
