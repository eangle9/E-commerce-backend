package cartservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type cartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) service.CartService {
	return &cartService{
		cartRepo: cartRepo,
	}
}

func (c cartService) AddToCart(request request.CartRequest, userId uint) response.Response {
	cartResponse, err := c.cartRepo.InsertCartItem(request, userId)
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         cartResponse,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "successfully added product item to shopping cart",
	}

	return response
}
