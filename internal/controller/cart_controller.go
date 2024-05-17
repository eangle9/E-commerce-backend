package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type cartController struct {
	engine      *router.Router
	cartService service.CartService
}

func NewCartController(engine *router.Router, cartService service.CartService) *cartController {
	return &cartController{
		engine:      engine,
		cartService: cartService,
	}
}

func (cart *cartController) InitCartRouter() {
	r := cart.engine
	api := r.Group("/cart")

	api.POST("/add", cart.addToCartHandler)
}

func (cart cartController) addToCartHandler(c *gin.Context) {
	var request request.CartRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.Set("error", errorResponse)
		return
	}

	validate := c.MustGet("validator").(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: customizer8.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := cart.cartService.AddToCart(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
