package productitemservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
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
	productItem := dto.ProductItem{
		ProductID:  request.ProductID,
		ColorID:    request.ColorID,
		SizeID:     request.SizeID,
		Price:      request.Price,
		Discount:   request.Discount,
		QtyInStock: request.QtyInStock,
	}

	id, image_url, err := p.itemRepo.InsertProductItem(request)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	productItem.ID = *id
	productItem.ImageUrl = image_url

	response := response.Response{
		Data:         productItem,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, product item created successfully!",
	}

	return response
}
