package controlers

import (
	"context"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	auth "github.com/NesterovYehor/textnest/services/auth_service/api"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
)

type AuthControler struct {
	userSrv  *services.UserService
	tokenSrv *services.JwtService
	log      *jsonlog.Logger
	auth.UnimplementedAuthServiceServer
}

func NewAuthControler(log *jsonlog.Logger, userService *services.UserService, tokenSrv *services.JwtService) *AuthControler {
	return &AuthControler{
		log:      log,
		userSrv:  userService,
		tokenSrv: tokenSrv,
	}
}

func (controler *AuthControler) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	err := controler.userSrv.CreateNewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &auth.CreateUserResponse{}, nil
}

func (controler *AuthControler) AuthenticateUser(ctx context.Context, req *auth.AuthenticateUserRequest) (*auth.AuthenticateUserResponse, error) {
	userId, err := controler.userSrv.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := controler.tokenSrv.GenerateAccessTocken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := controler.tokenSrv.GenerateRefreshTocken(userId)
	if err != nil {
		return nil, err
	}

	return &auth.AuthenticateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
