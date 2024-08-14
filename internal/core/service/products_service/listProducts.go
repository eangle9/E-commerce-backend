package productsservice

import (
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type productsService struct {
	productsRepo  repository.GetProducts
	serviceLogger *zap.Logger
}

func NewProductsService(productsRepo repository.GetProducts, serviceLogger *zap.Logger) service.GetProductService {
	return &productsService{
		productsRepo:  productsRepo,
		serviceLogger: serviceLogger,
	}
}

func (p productsService) GetAllProducts(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error) {
	page := request.Page
	perPage := request.PerPage

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 30
	}

	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed pagination query validation").WithProperty(entity.StatusCode, 400)
		p.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GetAllProducts"),
			zap.String("requestID", requestID),
			zap.Any("paginationData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	offset := (page - 1) * perPage

	productList, err := p.productsRepo.ListAllProducts(ctx, offset, perPage, requestID)
	if err != nil {
		return response.Response{}, err
	}

	data := response.Data{
		MetaData: response.PaginationQuery{
			Page:    page,
			PerPage: perPage,
		},
		Data: productList,
	}

	response := response.Response{
		Data:       data,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of reviews",
	}

	return response, nil
}
