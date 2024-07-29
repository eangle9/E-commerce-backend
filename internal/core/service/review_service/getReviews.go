package reviewservice

import (
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (r reviewService) GetReviews(ctx context.Context, paginationQuery request.PaginationQuery, requestID string) (response.Response, error) {
	page := paginationQuery.Page
	perPage := paginationQuery.PerPage

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 5
	}

	if err := paginationQuery.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed pagination query validation").WithProperty(entity.StatusCode, 400)
		r.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GetReviews"),
			zap.String("requestID", requestID),
			zap.Any("paginationData", paginationQuery),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	offset := (page - 1) * perPage
	reviews, err := r.reviewRepo.ListReviews(ctx, offset, perPage, requestID)

	if err != nil {
		return response.Response{}, err
	}

	data := response.Data{
		MetaData: response.PaginationQuery{
			Page:    page,
			PerPage: perPage,
		},
		Data: reviews,
	}

	response := response.Response{
		Data:       data,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of reviews",
	}

	return response, nil

}
