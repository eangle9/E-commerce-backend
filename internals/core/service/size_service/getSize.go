package sizeservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (s sizeService) GetSize(ctx context.Context, id int, requestID string) (response.Response, error) {
	size, err := s.sizeRepo.GetSizeById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       size,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get size with size_id '%d'", id),
	}

	return response, nil
}
