package colorservice

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (c colorService) GetColors() response.Response {
	colors, err := c.colorRepo.ListColors()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       colors,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of colors",
	}

	return response
}
