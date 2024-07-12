package service

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (u userService) DeleteUser(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := u.userRepo.DeleteUserById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("user with user_id %d deleted successfully", id),
	}
	return response, nil
}
