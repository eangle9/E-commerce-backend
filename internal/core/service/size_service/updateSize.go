package sizeservice

import (
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (s sizeService) UpdateSize(ctx context.Context, id int, size request.UpdateSize, requestID string) (response.Response, error) {
	updatedSize, err := s.sizeRepo.EditSizeById(ctx, id, size, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       updatedSize,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated the size with size_id '%d'", id),
	}

	return response, nil
}
