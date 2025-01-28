package validation

import (
	"errors"
	"regexp"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
)

// EmailRX is a regex to validate email addresses.
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// ValidateEmail validates an email address.
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is not provided")
	}
	if !validator.Match(email, EmailRX) {
		return errors.New("email is not valid")
	}
	return nil
}

// ValidatePasswordPlaintext validates a plaintext password.
func ValidatePasswordPlaintext(password string) error {
	if len([]rune(password)) < 8 {
		return errors.New("password is shorter than 8 characters")
	}
	if len([]rune(password)) > 72 {
		return errors.New("password is longer than 72 characters")
	}
	return nil
}

// ValidateUser validates a User model.
func ValidateUser(user *models.User) error {
	if user.Name == "" {
		return errors.New("name is not provided")
	}
	if len([]rune(user.Name)) > 100 {
		return errors.New("name is longer than 100 characters")
	}

	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	if user.Password.Plaintext != nil {
		return ValidatePasswordPlaintext(*user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		return errors.New("password hash is missing")
	}

	return nil
}
