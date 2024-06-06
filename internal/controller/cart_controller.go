package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"fmt"
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
	protectedMiddleware := middleware.ProtectedMiddleware
	r := cart.engine
	api := r.Group("/cart")

	api.POST("/add", protectedMiddleware(), cart.addToCartHandler)
}

// addToCartHandler godoc
// @Summary		    AddToCart
// @Description	    Add product to shopping cart
// @Tags			cart
// @ID				add-to-cart
// @Accept			json
// @Produce		    json
// @Security		JWT
// @Param			cart_item	body		request.CartRequest	true	"Cart item data"
// @Success		    200			{object}	response.Response
// @Router			/cart/add [post]
func (cart cartController) addToCartHandler(c *gin.Context) {
	var request request.CartRequest
	idAny, exist := c.Get("userId")
	if !exist {
		errorResponse := response.Response{
			Status:       http.StatusUnauthorized,
			ErrorType:    errorcode.Unauthorized,
			ErrorMessage: "unable to get userId from the context",
		}
		c.Set("error", errorResponse)
		return
	}
	userId := idAny.(uint)

	roleAny, exist := c.Get("role")
	if !exist {
		errorResponse := response.Response{
			Status:       http.StatusUnauthorized,
			ErrorType:    errorcode.Unauthorized,
			ErrorMessage: "unable to get role from the context",
		}
		c.Set("error", errorResponse)
		return
	}
	role := roleAny.(string)
	fmt.Println("role:", role)

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

	resp := cart.cartService.AddToCart(request, userId)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
