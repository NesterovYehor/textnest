// services/upload_service/internal/validation/metadata_validator.go
package validation

import (
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
)

// ValidateMetaData performs validation checks on metadata
func ValidateMetaData(metadata *pb.UploadPasteRequest) *validator.Validator {
	v := validator.New()
	v.Check(len([]rune(metadata.Key)) == 8, "key", "Key should be 8 characters long")
	v.Check(metadata.ExpirationDate.AsTime().After(time.Now()), "expiration_date", "Expiration date must be in the future")
	return v
}
