package controlers

import (
	"context"
	"errors"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	auth "github.com/NesterovYehor/textnest/services/auth_service/api"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	userSrv  *services.UserService
	tokenSrv *services.JwtService
	log      *jsonlog.Logger
	auth.UnimplementedAuthServiceServer
}

func NewAuthControler(log *jsonlog.Logger, userService *services.UserService, tokenSrv *services.JwtService) *AuthController {
	return &AuthController{
		log:      log,
		userSrv:  userService,
		tokenSrv: tokenSrv,
	}
}

func (ctr *AuthController) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	err := ctr.userSrv.CreateNewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &auth.CreateUserResponse{}, nil
}

func (ctr *AuthController) AuthenticateUser(ctx context.Context, req *auth.AuthenticateUserRequest) (*auth.AuthenticateUserResponse, error) {
	userId, err := ctr.userSrv.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, expiresAt, err := ctr.tokenSrv.GenerateAccessToken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExpiresAt, err := ctr.tokenSrv.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	return &auth.AuthenticateUserResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        timestamppb.New(expiresAt),
		RefreshExpiresAt: timestamppb.New(refreshExpiresAt),
	}, nil
}

func (ctr *AuthController) AuthorizeUser(ctx context.Context, req *auth.AuthorizeUserRequest) (*auth.AuthorizeUserResponse, error) {
	userId, err := ctr.tokenSrv.ExtractUserID(req.Tocken, "access")
	if err != nil {
		return nil, err
	}

	exist, err := ctr.userSrv.UserExists(userId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("User with this user id does not exist")
	}
	return &auth.AuthorizeUserResponse{
		UserId: userId,
	}, nil
}

func (ctr *AuthController) RefreshTokens(ctx context.Context, req *auth.RefreshTokensRequest) (*auth.RefreshTokensResponse, error) {
	userId, err := ctr.tokenSrv.ExtractUserID(req.Tocken, "refresh")
	if err != nil {
		return nil, status.Error(codes.NotFound, "User related to this token was not found")
	}

	accessToken, expiresAt, err := ctr.tokenSrv.GenerateAccessToken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExpiresAt, err := ctr.tokenSrv.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}
	return &auth.RefreshTokensResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        timestamppb.New(expiresAt),
		RefreshExpiresAt: timestamppb.New(refreshExpiresAt),
	}, nil
}
