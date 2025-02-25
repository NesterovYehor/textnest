package controllers

import (
	"context"
	"fmt"
	"sync"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	auth "github.com/NesterovYehor/textnest/services/auth_service/api"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/mailer"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	accessType   = "access"
	refreshType  = "refresh"
	activateType = "activate"
)

type AuthController struct {
	userSrv  *services.UserService
	tokenSrv *services.TokenService
	log      *jsonlog.Logger
	mailer   *mailer.Mailer
	auth.UnimplementedAuthServiceServer
}

func NewAuthController(log *jsonlog.Logger, userService *services.UserService, tokenSrv *services.TokenService, mailer *mailer.Mailer) *AuthController {
	return &AuthController{
		log:      log,
		userSrv:  userService,
		tokenSrv: tokenSrv,
		mailer:   mailer,
	}
}

func (ctr *AuthController) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	userID, err := ctr.userSrv.CreateNewUser(req.Name, req.Email, req.Password)
	if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create new user: %v", err))
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = ctr.mailer.Send(req.Email, "user_welcome.tmpl", map[string]string{
			"userID": userID,
		}); err != nil {
			ctr.log.PrintError(ctx, err, nil)
		}
	}()
	wg.Wait()

	return &auth.CreateUserResponse{}, nil
}

func (ctr *AuthController) ActivateUser(ctx context.Context, req *auth.ActivateUserRequest) (*auth.ActivateUserResponse, error) {
	// Attempt to activate the user
	if err := ctr.userSrv.ActivateUser(req.UserId); err != nil {
		return nil, status.Error(status.Code(err), "User activation failed. Please verify the user ID and try again.")
	}

	// Successfully activated the user
	return &auth.ActivateUserResponse{Message: "User has been successfully activated."}, nil
}

func (ctr *AuthController) AuthenticateUser(ctx context.Context, req *auth.AuthenticateUserRequest) (*auth.AuthenticateUserResponse, error) {
	userId, err := ctr.userSrv.AuthenticateUserByEmail(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Authentication failed. Check your credentials and try again.")
	}

	accessToken, expiresAt, err := ctr.tokenSrv.GenerateJWTToken(userId, accessType)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate access token. Please try again.")
	}

	refreshToken, refreshExpiresAt, err := ctr.tokenSrv.GenerateJWTToken(userId, refreshType)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate refresh token. Please try again.")
	}

	return &auth.AuthenticateUserResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        timestamppb.New(expiresAt),
		RefreshExpiresAt: timestamppb.New(refreshExpiresAt),
	}, nil
}

func (ctr *AuthController) AuthorizeUser(ctx context.Context, req *auth.AuthorizeUserRequest) (*auth.AuthorizeUserResponse, error) {
	userId, err := ctr.tokenSrv.ExtractUserID(req.Tocken, accessType)
	if err != nil {
		return nil, status.Error(status.Code(err), "Failed to extract user ID from the token. Please check the token and try again.")
	}

	exist, err := ctr.userSrv.UserExists(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check user existence. Please try again.")
	}
	if !exist {
		return nil, status.Error(codes.NotFound, "User with this ID does not exist.")
	}

	return &auth.AuthorizeUserResponse{
		UserId: userId,
	}, nil
}

func (ctr *AuthController) RefreshTokens(ctx context.Context, req *auth.RefreshTokensRequest) (*auth.RefreshTokensResponse, error) {
	userId, err := ctr.tokenSrv.ExtractUserID(req.Tocken, refreshType)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User related to this token was not found. Please ensure the token is valid.")
	}

	accessToken, expiresAt, err := ctr.tokenSrv.GenerateJWTToken(userId, accessType)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate access token. Please try again.")
	}

	refreshToken, refreshExpiresAt, err := ctr.tokenSrv.GenerateJWTToken(userId, refreshType)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate refresh token. Please try again.")
	}

	return &auth.RefreshTokensResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        timestamppb.New(expiresAt),
		RefreshExpiresAt: timestamppb.New(refreshExpiresAt),
	}, nil
}

func (ctr *AuthController) SendPasswordResetToken(ctx context.Context, req *auth.SendPasswordResetTokenRequest) (*auth.SendPasswordResetTokenResponse, error) {
	userID, err := ctr.userSrv.ValidateUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "User with this email does not exist.")
	}
	token, err := ctr.tokenSrv.CreateResetToken(userID)
	if err != nil {
		ctr.log.PrintError(ctx, err, nil)
		return nil, status.Error(codes.Internal, "Failed to generate reset token.")
	}
	go func() {
		if err := ctr.mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
			"passwordResetToken": token,
		}); err != nil {
			ctr.log.PrintError(ctx, err, nil)
		}
	}()
	return &auth.SendPasswordResetTokenResponse{Message: "Email with instractions for reseting password was send"}, nil
}

func (ctr *AuthController) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	var wg sync.WaitGroup
	if err := ctr.tokenSrv.ValidateResetToken(req.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Token for reseting password is invalid: %v", err))
	}
	userID, err := ctr.userSrv.ResetPassword(req.Password, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to reset password: %v", err))
	}
	wg.Add(1)
	go func() {
		if err := ctr.tokenSrv.DeleteAllForUser(*userID); err != nil {
			ctr.log.PrintError(ctx, err, nil)
		}
	}()
	wg.Wait()
	return &auth.ResetPasswordResponse{Message: "Password is renewed"}, nil
}
