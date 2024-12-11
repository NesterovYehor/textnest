package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConnectionIdle = 5 * time.Minute
	maxConnectionAge  = 5 * time.Minute
	gRPCTimeout       = 15 * time.Second
	gRPCTime          = 10 * time.Second
)

type GrpcConfig struct {
	Port string `yaml:"port"`
}

type GrpcServer struct {
	Grpc   *grpc.Server
	Config *GrpcConfig
}

func NewGrpcServer(cfg *GrpcConfig) *GrpcServer {
	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: maxConnectionIdle,
		Timeout:           gRPCTimeout,
		MaxConnectionAge:  maxConnectionAge,
		Time:              gRPCTime,
	}))

	return &GrpcServer{
		Grpc:   srv,
		Config: cfg,
	}
}

func (srv *GrpcServer) RunGrpcServer(ctx context.Context) error {
	address := fmt.Sprintf(":%s", srv.Config.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start gRPC server: %v", err)
	}

	fmt.Printf("gRPC server is listening on: %s\n", address)

	// Run server in a separate goroutine
	go func() {
		if err := srv.Grpc.Serve(listen); err != nil {
			fmt.Printf("[RunGrpcServer] gRPC server serve error: %v\n", err)
		}
	}()

	// Wait for context cancellation to stop the server
	<-ctx.Done()
	fmt.Printf("Shutting down gRPC server on PORT: %s\n", srv.Config.Port)
	srv.shutdown()
	fmt.Println("gRPC server exited properly")
	return nil
}

func (srv *GrpcServer) shutdown() {
	srv.Grpc.Stop()
	srv.Grpc.GracefulStop()
}
