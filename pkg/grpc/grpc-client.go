package grpc

import (
	"log"

	"google.golang.org/grpc"
)

func NewGrpcClient(addr string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Printf("Could not connect to gRPC server: %v", err)
		return nil, err
	}
	return conn, nil
}
