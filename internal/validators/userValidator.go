package validators

import (
	"errors"
	"regexp"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

var (
	ErrMissingFields    = errors.New("required fields missing")
	ErrNameStringLength = errors.New("fullname must be at least 4 letters longs")
	ErrPasswordLength   = errors.New("password must be at least 8 characters")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordMatch    = errors.New("passwords do not match")
)

func ValidateNewUserFields(newUser models.CreateUserInput) error {
	if len(newUser.Email) == 0 || len(newUser.Fullname) == 0 || len(newUser.Password) == 0 {
		return ErrMissingFields
	}
	if len(newUser.Fullname) < 4 {
		return ErrNameStringLength
	}

	if len(newUser.Password) < 8 {
		return ErrPasswordLength
	}
	if !isValidEmail(newUser.Email) {
		return ErrInvalidEmail
	}
	if newUser.Password != newUser.ConfirmPassword {
		return ErrPasswordMatch
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateUpdateUserFields(updateUser models.UpdateUserInput) error {
	if updateUser.Fullname != nil && len(*updateUser.Fullname) < 5 {
		return errors.New("fullname must be at least 3 characters long")
	}

	return nil
}
