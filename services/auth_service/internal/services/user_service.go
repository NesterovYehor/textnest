package services

import (
	"errors"
	"fmt"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/validation"
)

type UserService struct {
	model *models.UserModel
}

// NewUserService initializes a new UserService.
func NewUserService(model *models.UserModel) *UserService {
	return &UserService{model: model}
}

// CreateNewUser creates a new user and inserts it into the database.
func (srv *UserService) CreateNewUser(name, email, password string) error {
	newUser := &models.User{
		Name:      name,
		Email:     email,
		Activated: true,
	}

	if err := newUser.Password.Set(password); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := validation.ValidateUser(newUser); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	if err := srv.model.Insert(newUser); err != nil {
		if err == models.ErrDuplicateEmail {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("failed to insert new user into the database: %w", err)
	}

	return nil
}

func (srv *UserService) GetUserByEmail(email, passwordPlaintext string) (int64, error) {
	user, err := srv.model.GetByEmail(email)
	if err != nil {
		return 0, err
	}
	matches, err := user.Password.Match(passwordPlaintext)
	if err != nil {
		return 0, err
	}
	if !matches {
		return 0, errors.New("password is incorrect")
	}
	return user.ID, nil
}

func (srv *UserService) UserExists(userId int64) (bool, error) {
	return srv.model.UserExists(userId)
}
