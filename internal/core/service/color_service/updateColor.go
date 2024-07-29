package colorservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
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
