package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidUUID    = errors.New("invalid UUID")
	ErrDatabaseError  = errors.New("database error")
	ErrUpdateFailed   = errors.New("update failed")
	ErrInsertFailed   = errors.New("insert failed")
	ErrSelectFailed   = errors.New("select failed")
	AnonymousUser     = &User{}
)

type Model struct {
	User  *UserModel
	Token *TokenModel
}

func New(pool *pgxpool.Pool) *Model {
	return &Model{
		User:  NewUserModel(pool),
		Token: NewTokenModel(pool),
	}
}
