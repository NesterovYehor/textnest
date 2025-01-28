package services

import (
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
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

func (srv *JwtService) GenerateAccessTocken(userId int64) (string, error) {
	return srv.generateToken(userId, "access", srv.JwtConfig.AccessExpiry, []byte(srv.JwtConfig.AccessSecret))
}

func (srv *JwtService) GenerateRefreshTocken(userId int64) (string, error) {
	return srv.generateToken(userId, "refresh", srv.JwtConfig.RefreshExpiry, []byte(srv.JwtConfig.RefreshSecret))
}

func (srv *JwtService) generateToken(
	userId int64,
	tokenType string, // "access" or "refresh"
	expiry time.Duration, // e.g., 15m or 7d
	secret []byte, // secret key for signing
) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"type":    tokenType,
		"exp":     time.Now().Add(expiry),
	}
	token := jwt.NewWithClaims(srv.JwtConfig.SigningMethod, claims)
	return token.SignedString(secret) // Sign with the correct secret
}
