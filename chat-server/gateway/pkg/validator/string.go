package validator

import (
	"errors"
	"regexp"
)

func ValidateLogin(str string) error {
	if len(str) < 4 || len(str) > 32 {
		return errors.New("login string length must be between 4 and 32")
	}

	regex := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !regex.MatchString(str) {
		return errors.New("login string contains invalid characters")
	}

	return nil
}

func ValidatePassword(str string) error {
	if len(str) < 8 {
		return errors.New("password string length must be between 8 and 16")
	}
	return nil
}
