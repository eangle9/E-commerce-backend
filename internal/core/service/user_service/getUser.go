package service

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (u userService) GetUser(id int) (response.Response, error) {
	user, err := u.userRepo.GetUserById(id)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       user,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get the user with userID %d successfully", id),
	}
	return response, nil
}
