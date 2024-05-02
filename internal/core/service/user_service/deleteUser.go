package service

import (
	"Eccomerce-website/internal/core/model/response"
)

func (u userService) DeleteUser(id int) response.Response {
	resp, status, errType, err := u.userRepo.DeleteUserById(id)
	if err != nil {
		response := response.Response{
			Status:       status,
			ErrorType:    errType,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Status:       status,
		ErrorType:    errType,
		ErrorMessage: resp,
	}
	return response
}
