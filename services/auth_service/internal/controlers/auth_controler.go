package controlers

import (
	"context"
	"errors"

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

func (ctr *AuthControler) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	err := ctr.userSrv.CreateNewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &auth.CreateUserResponse{}, nil
}

func (ctr *AuthControler) AuthenticateUser(ctx context.Context, req *auth.AuthenticateUserRequest) (*auth.AuthenticateUserResponse, error) {
	userId, err := ctr.userSrv.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := ctr.tokenSrv.GenerateAccessTocken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := ctr.tokenSrv.GenerateRefreshTocken(userId)
	if err != nil {
		return nil, err
	}

	return &auth.AuthenticateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ctr *AuthControler) AuthorizeUser(ctx context.Context, req *auth.AuthorizeUserRequest) (*auth.AuthorizeUserResponse, error) {
	userId, err := ctr.tokenSrv.ExtractUserID(req.Tocken, "access")
	if err != nil {
		return nil, err
	}

	exist, err := ctr.userSrv.UserExists(*userId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("User with this user id does not exist")
	}
	return &auth.AuthorizeUserResponse{
		UserId: *userId,
	}, nil
}
