package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcConfig struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Development bool   `mapstructure:"development"`
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
	listen, err := net.Listen("tcp", srv.Config.Port)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("shutting down grpc PORT: {%s}", srv.Config.Port)
				srv.shutdown()
				fmt.Println("grpc exited properly")
				return
			}
		}
	}()

	fmt.Println("grpc server is listening on port: %s", srv.Config.Port)

	err = srv.Grpc.Serve(listen)
	if err != nil {
		fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err)
	}

	return err
}

func (srv *GrpcServer) shutdown() {
	srv.Grpc.Stop()
	srv.Grpc.GracefulStop()
}
