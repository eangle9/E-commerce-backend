package service

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
)

func (u userService) DeleteUser(id int) response.Response {
	resp, status, err := u.userRepo.DeleteUserById(id)
	if err != nil {
		response := response.Response{
			Status:       status,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Status:       status,
		ErrorType:    errorcode.Success,
		ErrorMessage: resp,
	}
	return response
}
