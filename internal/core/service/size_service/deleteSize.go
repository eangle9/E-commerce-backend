package sizeservice

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (s sizeService) DeleteSize(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := s.sizeRepo.DeleteSizeById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("size with size_id %d deleted successfully", id),
	}

	return response, nil
}
