package service

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (u userService) UpdateUser(id int, user utils.UpdateUser) response.Response {
	updateUser, err := u.userRepo.EditUserById(id, user)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return errorResponse
	}
	errorMessage := fmt.Sprintf("you have successfully updated the user with userId %d", id)
	response := response.Response{
		Data:         updateUser,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}
	return response
}
