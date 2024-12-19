package colorservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (c colorService) DeleteColor(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := c.colorRepo.DeleteColorById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("color with color_id %d deleted successfully", id),
	}

	return response, nil
}
