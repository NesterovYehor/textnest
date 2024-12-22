package validation

import "github.com/NesterovYehor/TextNest/pkg/validator"

func ValidateKey(key string) *validator.Validator {
	v := validator.New()
	v.Check(len([]rune(key)) != 8, "key", "Key must be 8 chars lenth")
	return v
}
