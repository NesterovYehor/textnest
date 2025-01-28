package grpc_clients

import (
	"context"
	"time"

	auth "github.com/NesterovYehor/TextNest/services/api_service/api/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

func NewAuthClient(target string) (*AuthClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		client: auth.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) SignUp(name, email, password string) (*auth.CreateUserResponse, error) {
	req := &auth.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.CreateUser(ctx, req)
}

func (c *AuthClient) LogIn(email, password string) (*auth.AuthenticateUserResponse, error) {
	req := auth.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := c.client.AuthenticateUser(ctx, &req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
