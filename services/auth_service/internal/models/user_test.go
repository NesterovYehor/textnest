package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword_Set(t *testing.T) {
	tests := []struct {
		name        string
		plaintext   string
		expectError bool
	}{
		{"Valid password", "securepassword", false},
		{"Empty password", "", false},                         // bcrypt allows empty passwords
		{"Long password", string(make([]byte, 1000)), true}, // Test edge case for long input
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p password

			err := p.Set(tt.plaintext)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p.Hash)
				assert.Equal(t, tt.plaintext, p.Plaintext)

				// Optional: Verify the hash matches the plaintext
				err = bcrypt.CompareHashAndPassword(p.Hash, []byte(tt.plaintext))
				assert.NoError(t, err, "Hash does not match plaintext")
			}
		})
	}
}
