package sizeservice

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (s sizeService) GetSizes() response.Response {
	sizes, err := s.sizeRepo.ListSizes()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       sizes,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of sizes",
	}

	return response
}
