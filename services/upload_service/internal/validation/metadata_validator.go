// services/upload_service/internal/validation/metadata_validator.go
package validation

import (
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
)

// ValidateMetaData performs validation checks on metadata
func ValidateMetaData(metadata *models.MetaData) *validator.Validator {
	v := validator.New()
	v.Check(metadata.UserId != "", "userId", "User id must be provided")
	v.Check(len([]rune(metadata.Key)) == 8, "key", "Key should be 8 characters long")
	v.Check(!metadata.CreatedAt.After(time.Now()), "created_at", "Paste creation date cannot be in the future")
	v.Check(metadata.ExpirationDate.After(time.Now()), "expiration_date", "Expiration date must be in the future")
	return v
}
