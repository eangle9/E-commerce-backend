package colorservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (c colorService) GetColor(ctx context.Context, id int, requestID string) (response.Response, error) {
	color, err := c.colorRepo.GetColorById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	errorMessage := fmt.Sprintf("you have get a color with color_id %d", id)
	response := response.Response{
		Data:       color,
		StatusCode: http.StatusOK,
		Message:    errorMessage,
	}

	return response, nil
}
