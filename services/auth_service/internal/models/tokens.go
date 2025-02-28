package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Token struct {
	Hash   string
	UserID uuid.UUID
	Expiry time.Time
}

type TokenModel struct {
	pool *pgxpool.Pool
}

func NewTokenModel(pool *pgxpool.Pool) *TokenModel {
	return &TokenModel{pool: pool}
}

func (m *TokenModel) Insert(token *Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	query := `
        INSERT INTO tokens (hash, user_id, expiry) VALUES ($1, $2, $3)
    `
	args := []any{
		token.Hash,
		token.UserID,
		token.Expiry,
	}
	_, err := m.pool.Exec(ctx, query, args...)
	return err
}

func (m *TokenModel) DeleteAllForUser(userID uuid.UUID) error {
	query := `
        DELETE FROM tokens
        WHERE user_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.pool.Exec(ctx, query, userID)
	return err
}

func (m *TokenModel) GetToken(tokenHash string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := `
        SELECT user_id, expiry FROM tokens WHERE hash = $1
    `
	var token Token
	err := m.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.UserID,
		&token.Expiry,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	token.Hash = tokenHash
	return &token, nil
}
