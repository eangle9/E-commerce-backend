package validationdata

import (
	"Eccomerce-website/internal/core/entity"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
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

func HasPassword(password string, serviceLogger *zap.Logger, requestID string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "hashing error").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("hash password error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "HasPassword"),
			zap.String("requestID", requestID),
			zap.String("password_length", fmt.Sprintf("%d", len(password))),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return "", errorResponse
	}
	hashPassword := string(hash)
	return hashPassword, nil
}

func MatchPassword(hashPassword string, password string, dbLogger *zap.Logger, requestID string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		errorResponse := entity.InvalidCredentials.Wrap(err, "invalid password").WithProperty(entity.StatusCode, 401)
		dbLogger.Error("password mismatch",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "MatchPassword"),
			zap.String("requestID", requestID),
			zap.String("password_length", fmt.Sprintf("%d", len(password))),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return false, errorResponse
	}

	return true, nil
}
