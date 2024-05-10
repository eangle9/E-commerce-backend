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

type productController struct {
	engine         *router.Router
	productService service.ProductService
}

func NewProductController(engine *router.Router, productService service.ProductService) *productController {
	return &productController{
		engine:         engine,
		productService: productService,
	}
}

func (p *productController) InitProductRouter() {
	r := p.engine
	api := r.Group("/product")

	api.POST("/create", p.createProductHandler)
}

func (p productController) createProductHandler(c *gin.Context) {
	var request request.ProductRequest

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
			ErrorMessage: customizer6.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productService.CreateProduct(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
