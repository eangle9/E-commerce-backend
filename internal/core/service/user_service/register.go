package service

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	PhoneRe = regexp.MustCompile(`^(\+251|251|0)?[79]\d{8}$`)
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type data struct {
	User dto.User `json:"user"`
}

func hasPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashPassword := string(hash)
	return hashPassword, err
}

func (u userService) SignUp(request request.SignUpRequest) response.Response {
	// phone number validation
	if isPhoneValid := PhoneRe.MatchString(request.PhoneNumber); !isPhoneValid {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "invalid phone number format.Please enter a valid phone number",
		}
		return errorResponse
	}

	if request.PhoneNumber[:1] == "0" {
		request.PhoneNumber = request.PhoneNumber[1:]
	}
	if request.PhoneNumber[:3] == "251" {
		request.PhoneNumber = request.PhoneNumber[3:]
	}
	if request.PhoneNumber[:4] != "+251" {
		request.PhoneNumber = "+251" + request.PhoneNumber
	}

	// password validation
	isValid := false
	isUpper := false
	isLower := false
	isDigit := false
	isSpecialChar := false

	specialChar := "!@#$%^&*+_-?></|"

	for _, char := range request.Password {
		if unicode.IsUpper(char) {
			isUpper = true
		}
		if unicode.IsLower(char) {
			isLower = true
		}
		if unicode.IsDigit(char) {
			isDigit = true
		}
		if strings.ContainsRune(specialChar, char) {
			isSpecialChar = true
		}
		if isUpper && isLower && isDigit && isSpecialChar {
			isValid = true
			break
		}
	}

	if !isValid {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "the password must contain at least one uppercase letter, lowercase letter, digit and special character",
		}
		return errorResponse
	}

	if len(request.Password) < 8 {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "password must be at least 8 characters long",
		}
		return errorResponse
	}

	// hash password
	hashPassword, err := hasPassword(request.Password)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    "HASHING_ERROR",
			ErrorMessage: err.Error(),
		}
		return errorResponse
	}

	request.Password = hashPassword

	user := dto.User{
		Username:      request.Username,
		Email:         request.Email,
		Password:      request.Password,
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		PhoneNumber:   request.PhoneNumber,
		Address:       request.Address,
		Role:          request.Role,
		EmailVerified: false,
		// CreatedAt:     time.Now(),
		// UpdatedAt:     time.Now(),
		// DeletedAt:     gorm.DeletedAt{},
	}

	id, err := u.userRepo.InsertUser(user)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return errorResponse
	}

	user.ID = id
	data := data{
		User: user,
	}
	response := response.Response{
		Data:         data,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, you have registered successfully!",
	}
	return response
}
