package service

import (
	validationdata "Eccomerce-website/internal/core/common/utils/validationData"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
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

func (u userService) SignUp(request request.SignUpRequest) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed register validation").WithProperty(entity.StatusCode, 400)
		return response.Response{}, errorResponse
	}

	// hash password
	hashPassword, err := validationdata.HasPassword(request.Password)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "hashingError").WithProperty(entity.StatusCode, 500)
		return response.Response{}, errorResponse
	}

	request.Password = hashPassword

	user := dto.User{
		Username:      request.Username,
		Email:         request.Email,
		Password:      request.Password,
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		PhoneNumber:   request.PhoneNumber.Number,
		EmailVerified: false,
	}

	id, err := u.userRepo.InsertUser(user)
	if err != nil {
		return response.Response{}, err
	}

	user.ID = id
	data := data{
		User: user,
	}
	response := response.Response{
		Data:       data,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, you have registered successfully!",
	}
	return response, nil
}
