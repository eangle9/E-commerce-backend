package colorservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (c colorService) UpdateColor(id int, color utils.UpdateColor) response.Response {
	updatedColor, err := c.colorRepo.EditColorById(id, color)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       updatedColor,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of colors",
	}

	return response
}
