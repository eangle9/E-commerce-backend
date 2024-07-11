package reviewservice

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (r reviewService) GetReviews() response.Response {
	reviews, err := r.reviewRepo.ListReviews()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       reviews,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of review",
	}

	return response
}
