package colorservice

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (c colorService) UpdateColor(ctx context.Context, id int, color utils.UpdateColor, requestID string) (response.Response, error) {
	updatedColor, err := c.colorRepo.EditColorById(ctx, id, color, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       updatedColor,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated the color with color_id '%d'", id),
	}

	return response, nil
}
