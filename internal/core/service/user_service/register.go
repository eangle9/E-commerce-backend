package service

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/common/utils/password"
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
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

func (u userService) SignUp(request request.SignUpRequest) response.Response {
	// phone number validation
	errorResponse, phoneNumber := utils.PhoneValidation(request.PhoneNumber)
	if errorResponse != nil {
		return *errorResponse
	}

	// password validation
	if err := utils.PasswordValidation(request.Password); err != nil {
		return *err
	}

	// hash password
	hashPassword, err := password.HasPassword(request.Password)
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
		Username:    request.Username,
		Email:       request.Email,
		Password:    request.Password,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: phoneNumber,
		// Address:       request.Address,
		// Role:          "user",
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
