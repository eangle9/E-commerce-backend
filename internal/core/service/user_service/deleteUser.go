package service

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (u userService) DeleteUser(id int) (response.Response, error) {
	err := u.userRepo.DeleteUserById(id)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("user with user_id %d deleted successfully", id),
	}
	return response, nil
}
