package productitemservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"fmt"
	"net/http"
)

type productItemService struct {
	itemRepo repository.ProductItemRepository
}

func NewProductItemService(itemRepo repository.ProductItemRepository) service.ProductItemService {
	return &productItemService{
		itemRepo: itemRepo,
	}
}

func (p productItemService) CreateProductItem(request request.ProductItemRequest) response.Response {
	fmt.Printf("request: %+v", request)
	productItem := dto.ProductItem{
		ProductID: request.ProductID,
		ColorID:   request.ColorID,
		// SizeID:     request.SizeID,
		Price:      request.Price,
		Discount:   request.Discount,
		QtyInStock: request.QtyInStock,
	}

	id, image_url, err := p.itemRepo.InsertProductItem(request)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		}
		return response
	}

	productItem.ID = *id
	productItem.ImageUrl = image_url

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, product item created successfully!",
	}

	return response
}
