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
	if expectedType != "access" && expectedType != "refresh" {
		return "", fmt.Errorf("invalid token type: %s", expectedType)
	}
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

func (srv *JwtService) GenerateAccessToken(userId string) (string, time.Time, error) {
	return srv.generateToken(userId, "access")
}

func (srv *JwtService) GenerateRefreshToken(userId string) (string, time.Time, error) {
	return srv.generateToken(userId, "refresh")
}

func (srv *JwtService) generateToken(
	userId string,
	tokenType string, // "access" or "refresh"
) (string, time.Time, error) {
	var expiry time.Duration
	var secret []byte
	switch tokenType {
	case "access":
		expiry = srv.JwtConfig.AccessExpiry
		secret = []byte(srv.JwtConfig.AccessSecret)
	case "refresh":
		expiry = srv.JwtConfig.RefreshExpiry
		secret = []byte(srv.JwtConfig.RefreshSecret)
	default:
		return "", time.Time{}, fmt.Errorf("invalid token type: %s", tokenType)
	}
	claims := jwt.MapClaims{
		"user_id": userId,
		"type":    tokenType,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(srv.JwtConfig.SigningMethod, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenStr, time.Now().Add(expiry), nil
}
