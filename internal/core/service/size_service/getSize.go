package sizeservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
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
				Status:       http.StatusNotFound,
				ErrorType:    errorcode.NotFoundError,
				ErrorMessage: err.Error(),
			}
			return response
		}
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         size,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have get size with size_id '%d'", id),
	}

	return response
}
