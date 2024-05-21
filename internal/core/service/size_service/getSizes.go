package sizeservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (s sizeService) GetSizes() response.Response {
	sizes, err := s.sizeRepo.ListSizes()
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         sizes,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get all list of sizes",
	}

	return response
}
