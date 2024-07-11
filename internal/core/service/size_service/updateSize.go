package sizeservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (s sizeService) UpdateSize(id int, size utils.UpdateSize) response.Response {
	updatedSize, err := s.sizeRepo.EditSizeById(id, size)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       updatedSize,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated the size with size_id '%d'", id),
	}

	return response
}
