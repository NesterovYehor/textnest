package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	model     *models.TokenModel
	jwtConfig *config.JwtConfig
}

func NewTokenService(jwtCfg *config.JwtConfig, model *models.TokenModel) *TokenService {
	return &TokenService{
		jwtConfig: jwtCfg,
		model:     model,
	}
}

func (srv *TokenService) ExtractUserID(token string, expectedType string) (string, error) {
	if expectedType != "access" && expectedType != "refresh" {
		return "", fmt.Errorf("invalid token type: %s", expectedType)
	}
	secret := ""
	if expectedType == "access" {
		secret = srv.jwtConfig.AccessSecret
	} else {
		secret = srv.jwtConfig.RefreshSecret
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

func (srv *TokenService) GenerateJWTToken(
	userId string,
	tokenType string, // "access" or "refresh"
) (string, time.Time, error) {
	var expiry time.Duration
	var secret []byte
	switch tokenType {
	case "access":
		expiry = srv.jwtConfig.AccessExpiry
		secret = []byte(srv.jwtConfig.AccessSecret)
	case "refresh":
		expiry = srv.jwtConfig.RefreshExpiry
		secret = []byte(srv.jwtConfig.RefreshSecret)
	case "activate":
		expiry = srv.jwtConfig.ActivateExpiry
		secret = []byte(srv.jwtConfig.ActivateSecret)
	default:
		return "", time.Time{}, fmt.Errorf("invalid token type: %s", tokenType)
	}
	claims := jwt.MapClaims{
		"user_id": userId,
		"type":    tokenType,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(srv.jwtConfig.SigningMethod, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenStr, time.Now().Add(expiry), nil
}

func (srv *TokenService) GenerateSecureToken(userID *uuid.UUID) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil
	}
	hash := sha256.Sum256(randomBytes)
	hashedToken := hex.EncodeToString(hash[:])

	token := models.Token{
		UserID: *userID,
		Expiry: time.Now().Add(24 * time.Hour),
		Hash:   hashedToken,
	}

	if err := srv.model.Insert(&token); err != nil {
		return "", err
	}
	return token.Hash, nil
}

func (srv *TokenService) ValidateResetToken(tokenHash string) error {
	token, err := srv.model.GetToken(tokenHash)
	if err != nil {
		return err
	}
	return validation.ValidatePasswordResetToken(token)
}

func (srv *TokenService) DeleteAllForUser(userID uuid.UUID) error {
	return srv.model.DeleteAllForUser(userID)
}
