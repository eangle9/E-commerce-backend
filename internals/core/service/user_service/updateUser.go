package service

import (
	validationdata "Eccomerce-website/internals/core/common/utils/validationData"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (u userService) UpdateUser(ctx context.Context, id int, request request.UpdateUser, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed updateUser data validation").WithProperty(entity.StatusCode, 400)
		u.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "UpdateUser"),
			zap.String("requestID", requestID),
			zap.Int("id", id),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	if request.PhoneNumber != "" {
		formatedPhone, err := validationdata.FormatPhoneNumber(request.PhoneNumber)
		if err != nil {
			return response.Response{}, err
		}
		request.PhoneNumber = formatedPhone
	}

	if request.Password != "" {
		hashPassword, err := validationdata.HasPassword(request.Password, u.serviceLogger, requestID)
		if err != nil {
			return response.Response{}, err
		}

		request.Password = hashPassword
	}

	updateUser, err := u.userRepo.EditUserById(ctx, id, request, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       updateUser,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated the user with userId %d", id),
	}
	return response, nil
}
