package service

import (
	validationdata "Eccomerce-website/internal/core/common/utils/validationData"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type userService struct {
	userRepo      repository.UserRepository
	serviceLogger *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, serviceLogger *zap.Logger) service.UserService {
	return &userService{
		userRepo:      userRepo,
		serviceLogger: serviceLogger,
	}
}

type data struct {
	User dto.User `json:"user"`
}

func (u userService) SignUp(ctx context.Context, request request.SignUpRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed register validation").WithProperty(entity.StatusCode, 400)
		u.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "SignUp"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	// hash password
	hashPassword, err := validationdata.HasPassword(request.Password, u.serviceLogger, requestID)
	if err != nil {
		return response.Response{}, err
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

	id, err := u.userRepo.InsertUser(ctx, user, requestID)
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
