package colorservice

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (c colorService) UpdateColor(id int, color utils.UpdateColor) response.Response {
	updatedColor, err := c.colorRepo.EditColorById(id, color)
	if err != nil {
		response := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have successfully updated the color with color_id '%d'", id)
	response := response.Response{
		Data:         updatedColor,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}

	return response
}
