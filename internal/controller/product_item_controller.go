package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type productItemController struct {
	engine             *router.Router
	productItemService service.ProductItemService
}

func NewProductItemController(engine *router.Router, productItemService service.ProductItemService) *productItemController {
	return &productItemController{
		engine:             engine,
		productItemService: productItemService,
	}
}

func (p *productItemController) InitProductItemRouter() {
	r := p.engine
	api := r.Group("/item")

	api.POST("/create", p.createProductItemHandler)
	api.GET("/list", p.getProductItemsHandler)
	api.GET("/:id", p.getProductItemHandler)
	api.PUT("/update/:id", p.updateProductItemHandler)
}

func (p productItemController) createProductItemHandler(c *gin.Context) {
	var request request.ProductItemRequest

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
			ErrorMessage: customizer7.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.CreateProductItem(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (p productItemController) getProductItemsHandler(c *gin.Context) {
	resp := p.productItemService.GetProductItems()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (p productItemController) getProductItemHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer value",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.GetProductItem(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (p productItemController) updateProductItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	var request utils.UpdateProductItem

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.Set("error", errorResponse)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer value",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.UpdateProductItem(id, request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
