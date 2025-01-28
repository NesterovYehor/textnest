package validation_test

import (
	"testing"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/validation"
	"github.com/stretchr/testify/assert"
)

const (
	validEmail    = "valid@email.com"
	invalidEmail  = "invalidemail.com"
	validPassword = "12345678"
)

func TestValidateEmail(t *testing.T) {
	// Test valid email
	err := validation.ValidateEmail(validEmail)
	assert.NoError(t, err, "Valid email should not return an error")

	// Test invalid email
	err = validation.ValidateEmail(invalidEmail)
	assert.EqualError(t, err, "email is not valid", "Invalid email should return the correct error message")

	// Test empty email
	err = validation.ValidateEmail("")
	assert.EqualError(t, err, "email is not provided", "Empty email should return the correct error message")
}

func TestValidatePasswordPlaintext(t *testing.T) {
	// Test valid password
	err := validation.ValidatePasswordPlaintext(validPassword)
	assert.NoError(t, err, "Valid password should not return an error")

	// Test short password
	err = validation.ValidatePasswordPlaintext("short")
	assert.EqualError(t, err, "password is shorter than 8 characters", "Short password should return the correct error message")

	// Test overly long password
	longPassword := string(make([]rune, 73))
	err = validation.ValidatePasswordPlaintext(longPassword)
	assert.EqualError(t, err, "password is longer than 72 characters", "Long password should return the correct error message")
}

func TestValidateUser(t *testing.T) {
	// Valid user
	user := &models.User{
		Name:  "Valid User",
		Email: validEmail,
	}
	user.Password.Set(validPassword)
	err := validation.ValidateUser(user)
	assert.NoError(t, err, "Valid user should not return an error")

	// User with no name
	user.Name = ""
	err = validation.ValidateUser(user)
	assert.EqualError(t, err, "name is not provided", "User with no name should return the correct error message")

	// User with overly long name
	user.Name = string(make([]rune, 101))
	err = validation.ValidateUser(user)
	assert.EqualError(t, err, "name is longer than 100 characters", "User with overly long name should return the correct error message")

	// User with invalid email
	user.Name = "Valid User"
	user.Email = invalidEmail
	err = validation.ValidateUser(user)
	assert.EqualError(t, err, "email is not valid", "User with invalid email should return the correct error message")

	// User with no password
	user.Email = validEmail
	user.Password.Plaintext = nil
	user.Password.Hash = nil
	err = validation.ValidateUser(user)
	assert.EqualError(t, err, "password hash is missing", "User with no password hash should return the correct error message")

	// User with valid hash but no plaintext password
	user.Password.Hash = []byte("valid-hash")
	err = validation.ValidateUser(user)
	assert.NoError(t, err, "User with valid hash but no plaintext password should not return an error")
}
