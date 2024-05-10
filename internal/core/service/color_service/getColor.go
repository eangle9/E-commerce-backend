package colorservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (c colorService) GetColor(id int) response.Response {
	color, err := c.colorRepo.GetColorById(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have get a color with color_id %d", id)
	response := response.Response{
		Data:         color,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}
	return response
}
