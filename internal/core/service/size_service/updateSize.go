package sizeservice

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (s sizeService) UpdateSize(id int, size utils.UpdateSize) response.Response {
	updatedSize, err := s.sizeRepo.EditSizeById(id, size)
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         updatedSize,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have successfully updated the size with size_id '%d'", id),
	}

	return response
}
