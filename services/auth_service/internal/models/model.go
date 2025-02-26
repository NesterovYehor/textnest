package models

import "github.com/NesterovYehor/textnest/services/auth_service/internal/database"

type Model struct {
	User  *UserModel
	Token *TokenModel
}

func New(db database.DB) *Model {
	return &Model{
		User:  NewUserModel(db.Conn),
		Token: NewTokenModel(db.Conn),
	}
}
