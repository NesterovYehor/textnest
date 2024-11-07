package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const maxConnectionIdle = 5 * time.Minute
const maxConnectionAge = 5 * time.Minute
const gRPCTimeout = 15 * time.Second
const gRPCTime = 10 * time.Second

type GrpcConfig struct {
	Port string `mapstrucTure:"port"`
	Host string `mapstructure:"host"`
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
	address := fmt.Sprintf("%s:%s", srv.Config.Host, srv.Config.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("shutting down grpc PORT: %s\n", srv.Config.Port)
				srv.shutdown()
				fmt.Println("grpc exited properly")
				return
			}
		}
	}()

	fmt.Printf("grpc server is listening on: %s\n", address)

	err = srv.Grpc.Serve(listen)
	if err != nil {
		fmt.Printf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err)
	}

	return err
}

func (srv *GrpcServer) shutdown() {
	srv.Grpc.Stop()
	srv.Grpc.GracefulStop()
}
