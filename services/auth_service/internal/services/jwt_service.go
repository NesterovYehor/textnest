package services

import (
	"fmt"
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/validation"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	JwtConfig *config.JwtConfig
}

func NewJwtService(jwtCfg *config.JwtConfig) *JwtService {
	return &JwtService{
		JwtConfig: jwtCfg,
	}
}

func (srv *JwtService) ExtractUserID(token string, expectedType string) (string, error) {
	secret := ""
	if expectedType == "access" {
		secret = srv.JwtConfig.AccessSecret
	} else {
		secret = srv.JwtConfig.RefreshSecret
	}
	// Validate the token
	parsedToken, err := validation.ValidateJwtToken(token, secret, expectedType)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unable to parse token claims")
	}

	// Extract user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found or invalid in token claims")
	}

	// Return user ID
	return userID, nil
}

func (srv *JwtService) GenerateAccessTocken(userId string) (string, error) {
	return srv.generateToken(userId, "access", srv.JwtConfig.AccessExpiry, []byte(srv.JwtConfig.AccessSecret))
}

func (srv *JwtService) GenerateRefreshTocken(userId string) (string, error) {
	return srv.generateToken(userId, "refresh", srv.JwtConfig.RefreshExpiry, []byte(srv.JwtConfig.RefreshSecret))
}

func (srv *JwtService) generateToken(
	userId string,
	tokenType string, // "access" or "refresh"
	expiry time.Duration, // e.g., 15m or 7d
	secret []byte, // secret key for signing
) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"type":    tokenType,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(srv.JwtConfig.SigningMethod, claims)
	return token.SignedString(secret) // Sign with the correct secret
}
