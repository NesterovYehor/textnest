package grpc_clients

import (
	"context"
	"fmt"
	"sync"
	"time"

	paste_download "github.com/NesterovYehor/TextNest/services/api_service/api/download_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	downloadByKeyReqPool = sync.Pool{
		New: func() any {
			return new(paste_download.DownloadByKeyRequest)
		},
	}
	downloadByUserIdReqPool = sync.Pool{
		New: func() any {
			return new(paste_download.DownloadByUserIdRequest)
		},
	}
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
	req := downloadByKeyReqPool.Get().(*paste_download.DownloadByKeyRequest)
	req.Key = key
	defer downloadByKeyReqPool.Put(req)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.client.DownloadByKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *DownloadClient) DownloadByUserId(userId string, limit, offset int32) (*paste_download.DownloadByUserIdResponse, error) {
	req := downloadByUserIdReqPool.Get().(*paste_download.DownloadByUserIdRequest)
	req.UserId = userId
	req.Limit = limit
	req.Offset = offset
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.DownloadByUserId(ctx, req)
	if err != nil {
		return nil, err
	}
    fmt.Println("Raw gRPC response: ", resp)

	return resp, nil
}
