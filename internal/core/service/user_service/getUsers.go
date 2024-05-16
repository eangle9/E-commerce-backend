package service

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (u userService) GetUsers() response.Response {
	users, err := u.userRepo.ListUsers()
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: "failed to get list of users",
		}
		return errorResponse
	}

	response := response.Response{
		Data:         users,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get list of users successfully!",
	}

	return response

}
