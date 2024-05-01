package service

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (u userService) GetUser(id int) response.Response {
	user, err := u.userRepo.GetUserById(id)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return errorResponse
	}
	errorMessage := fmt.Sprintf("you have get the user with userID %d successfully", id)
	response := response.Response{
		Data:         user,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}
	return response
}
