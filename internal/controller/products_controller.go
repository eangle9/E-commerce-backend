package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/port/service"

	"github.com/gin-gonic/gin"
)

type productsController struct {
	engine          *router.Router
	productsService service.GetProductService
}

func NewProductsController(engine *router.Router, productsService service.GetProductService) *productsController {
	return &productsController{
		engine:          engine,
		productsService: productsService,
	}
}

func (p *productsController) InitProductsRouter() {
	r := p.engine
	api := r.Group("/products")

	api.GET("/list", p.listProductsHandler)
}

// listProductsHandler godoc
//
//	@Summary		List of products
//	@Description	Retrieves a list of products
//	@Tags			products
//	@ID				list_of_products
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/products/list [get]
func (p productsController) listProductsHandler(c *gin.Context) {
	resp := p.productsService.GetAllProducts()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
