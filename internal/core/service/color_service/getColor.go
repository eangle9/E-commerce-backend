package colorservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (c colorService) GetColor(id int) response.Response {
	color, err := c.colorRepo.GetColorById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have get a color with color_id %d", id)
	response := response.Response{
		Data:       color,
		StatusCode: http.StatusOK,
		Message:    errorMessage,
	}
	return response
}
