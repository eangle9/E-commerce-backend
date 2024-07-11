package sizeservice

import (
	"Eccomerce-website/internal/core/model/response"
	"database/sql"
	"fmt"
	"net/http"
)

func (s sizeService) GetSize(id int) response.Response {
	size, err := s.sizeRepo.GetSizeById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := response.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			}
			return response
		}
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       size,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get size with size_id '%d'", id),
	}

	return response
}
