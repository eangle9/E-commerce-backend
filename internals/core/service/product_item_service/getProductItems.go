package productitemservice

import (
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (p productItemService) GetProductItems(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error) {
	page := request.Page
	perPage := request.PerPage

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 5
	}

	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed pagination query validation").WithProperty(entity.StatusCode, 400)
		p.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GetProductItems"),
			zap.String("requestID", requestID),
			zap.Any("paginationData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	offset := (page - 1) * perPage

	productItems, err := p.itemRepo.ListProductItems(ctx, offset, perPage, requestID)
	if err != nil {
		return response.Response{}, err
	}

	data := response.Data{
		MetaData: response.PaginationQuery{
			Page:    page,
			PerPage: perPage,
		},
		Data: productItems,
	}

	response := response.Response{
		Data:       data,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get single product item "),
	}

	return response, nil
}
