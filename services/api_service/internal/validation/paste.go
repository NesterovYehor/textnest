package validation

import (
	"errors"
	"time"
)

// PasteInput represents the input for a paste request.
type PasteInput struct {
	Title          string    `json:"title"`
	ExpirationDate time.Time `json:"expiration_date"`
}

// ValidatePasteInput validates the input for uploading a paste.
func ValidatePasteInput(input *PasteInput) error {
	if input.ExpirationDate.Before(time.Now()) {
		return errors.New("expiration date cannot be in the past")
	}
	return nil
}
