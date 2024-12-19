package validation

import (
	"errors"
	"time"
)

// PasteInput represents the input for a paste request.
type PasteInput struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Content        string    `json:"content"`
}

// ValidatePasteInput validates the input for uploading a paste.
func ValidatePasteInput(input *PasteInput) error {
	if input.Content == "" {
		return errors.New("content cannot be empty")
	}
	if input.ExpirationDate.Before(time.Now()) {
		return errors.New("expiration date cannot be in the past")
	}
	return nil
}
