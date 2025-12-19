package validators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

func ValidateNewUserFields(newUser models.CreateUserInput) []error {
	validationErrors := []error{}
	if len(newUser.Email) == 0 || len(newUser.Fullname) == 0 || len(newUser.Password) == 0 {
		validationErrors = append(validationErrors, errors.New("all fileds are required"))
	}
	if len(newUser.Fullname) < 5 {
		validationErrors = append(validationErrors, errors.New("fullname must be at least 3 characters long"))
	}

	if len(newUser.Password) < 8 {
		validationErrors = append(validationErrors, errors.New("password must be at least 6 characters long"))
	}
	if !isValidEmail(newUser.Email) {
		validationErrors = append(validationErrors, errors.New("invalid email format"))
	}
	fmt.Println("Validation Errors:", validationErrors)
	if len(validationErrors) > 0 {
		return validationErrors
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
