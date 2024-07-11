package cartservice

import (
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"database/sql"
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
		if err == sql.ErrNoRows {
			response := response.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			}
			return response
		} else {
			response := response.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
			return response
		}

	}

	response := response.Response{
		Data:       cartResponse,
		StatusCode: http.StatusOK,
		Message:    "successfully added product item to shopping cart",
	}

	return response
}
