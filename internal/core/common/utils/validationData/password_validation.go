package validationdata

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(value interface{}) error {
	password, exist := value.(string)
	if !exist {
		return validation.NewError("validation_password", "password doesn't exist")
	}

	if len(password) < 8 {
		return validation.NewError("validation_password", "password must be at least 8 characters long")
	}

	upper := regexp.MustCompile(`[A-Z]`)
	if !upper.MatchString(password) {
		return validation.NewError("validation_password", "password must contain at least one uppercase letter")
	}

	lower := regexp.MustCompile(`[a-z]`)
	if !lower.MatchString(password) {
		return validation.NewError("validation_password", "password must contain at least one lowercase letter")
	}

	digit := regexp.MustCompile(`[0-9]`)
	if !digit.MatchString(password) {
		return validation.NewError("validation_password", "password must contain at least one number")
	}

	special := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",\.<>\/?\\|]`)
	if !special.MatchString(password) {
		return validation.NewError("validation_password", "password must contain at least one special character")
	}

	return nil

}

func HasPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashPassword := string(hash)
	return hashPassword, err
}
