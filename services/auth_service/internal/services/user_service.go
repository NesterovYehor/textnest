package services

import (
	"errors"
	"fmt"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/validation"
	"github.com/google/uuid"
)

type UserService struct {
	model *models.UserModel
}

// NewUserService initializes a new UserService.
func NewUserService(model *models.UserModel) *UserService {
	return &UserService{model: model}
}

// CreateNewUser creates a new user and inserts it into the database.
func (srv *UserService) CreateNewUser(name, email, password string) (string, error) {
	newUser := &models.User{
		Name:  name,
		Email: email,
	}

	if err := newUser.Password.Set(password); err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	if err := validation.ValidateUser(newUser); err != nil {
		return "", fmt.Errorf("user validation failed: %w", err)
	}
	userID, err := srv.model.Insert(newUser)
	if err != nil {
		return "", srv.handleErr(err)
	}

	return userID, nil
}

func (srv *UserService) AuthenticateUserByEmail(email, passwordPlaintext string) (string, error) {
	user, err := srv.model.GetByEmail(email)
	if err != nil {
		return "", srv.handleErr(err)
	}

	matches, err := user.Password.Match(passwordPlaintext)
	if err != nil {
		return "", srv.handleErr(err)
	}

	if !matches {
		return "", errors.New("invalid credentials")
	}

	return user.ID.String(), nil
}

func (srv *UserService) ValidateUserByEmail(email string) (*uuid.UUID, error) {
	if err := validation.ValidateEmail(email); err != nil {
		return nil, err
	}
	user, err := srv.model.GetByEmail(email)
	if err != nil {
		return nil, srv.handleErr(err)
	}

	if !user.Activated {
		return nil, errors.New("User is not activated")
	}

	return &user.ID, nil
}

func (srv *UserService) UserExists(userId string) (bool, error) {
	res, err := srv.model.UserExists(userId)
	if err != nil {
		return false, srv.handleErr(err)
	}
	return res, nil
}

func (srv *UserService) ActivateUser(userID string) error {
	if err := srv.model.ActivateUser(userID); err != nil {
		return srv.handleErr(err)
	}
	return nil
}

func (srv *UserService) ResetPassword(plainText, token string) (*uuid.UUID, error) {
	if err := validation.ValidatePasswordPlaintext(plainText); err != nil {
		return nil, err
	}
	userID, err := srv.model.ResetPassword(plainText, token)
	if err != nil {
		return nil, srv.handleErr(err)
	}
	return userID, nil
}

func (srv *UserService) handleErr(err error) error {
	switch err {
	case models.ErrRecordNotFound:
		return fmt.Errorf("user not found: %w", err)
	case models.ErrDatabaseError:
		return fmt.Errorf("database error: %w", err)
	case models.ErrUpdateFailed:
		return fmt.Errorf("failed to update data: %w", err)
	case models.ErrInsertFailed:
		return fmt.Errorf("failed to insert data: %w", err)
	case models.ErrSelectFailed:
		return fmt.Errorf("failed to retrieve data: %w", err)
	case models.ErrDuplicateEmail:
		return fmt.Errorf("email already in use: %w", err)
	case models.ErrInvalidUUID:
		return fmt.Errorf("invalid UUID format: %w", err)
	default:
		return fmt.Errorf("unexpected error: %w", err)
	}
}
