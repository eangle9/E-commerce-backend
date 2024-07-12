package service

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (u userService) GetUser(ctx context.Context, id int, requestID string) (response.Response, error) {
	user, err := u.userRepo.GetUserById(ctx, id, requestID)
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
