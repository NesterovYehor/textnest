package validation

import (
	"errors"
)

func ValidateKey(key string) error {
	if len(key) != 8 {
		return errors.New("Key mast be 8 characters lenth")
	}
	return nil
}

func IsUserIdValid(userId string) error {
	if userId == "" {
		return errors.New("User Id is not provided")
	}
	return nil
}
