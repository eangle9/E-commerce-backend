package productitemservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type productItemService struct {
	itemRepo      repository.ProductItemRepository
	serviceLogger *zap.Logger
}

func NewProductItemService(itemRepo repository.ProductItemRepository, serviceLogger *zap.Logger) service.ProductItemService {
	return &productItemService{
		itemRepo:      itemRepo,
		serviceLogger: serviceLogger,
	}
}

func (p productItemService) CreateProductItem(ctx context.Context, request request.ProductItemRequest, requestID string) (response.Response, error) {
	if request.Discount.LessThan(request.Price) {
		err := errors.New("discount can't be less than product price")
		errorResponse := entity.AppInternalError.Wrap(err, "discount must be greater than or equal to price").WithProperty(entity.StatusCode, 500)
		return response.Response{}, errorResponse
	}

	productItem := dto.ProductItem{
		ProductID:  request.ProductID,
		ColorID:    request.ColorID,
		Price:      request.Price,
		Discount:   request.Discount,
		QtyInStock: request.QtyInStock,
	}

	id, image_url, err := p.itemRepo.InsertProductItem(ctx, request, requestID)
	if err != nil {
		return response.Response{}, err
	}

	productItem.ID = *id
	productItem.ImageUrl = image_url

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, product item created successfully!",
	}

	return response, nil
}
