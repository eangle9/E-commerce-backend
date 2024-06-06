package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"net/http"
	"strconv"

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
	protectedMiddleware := middleware.ProtectedMiddleware
	r := p.engine
	api := r.Group("/product")

	api.POST("/create", protectedMiddleware(), p.createProductHandler)
	api.GET("/list", p.listProductHandler)
	api.GET("/:id", p.getProductHandler)
	api.PUT("/update/:id", protectedMiddleware(), p.updateProductHandler)
	api.DELETE("/delete/:id", protectedMiddleware(), p.deleteProductHandler)
}

// createProductHandler godoc
//	@Summary		create product
//	@Description	Insert a new product
//	@Tags			product
//	@ID				create-product
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			product	body		request.ProductRequest	true	"Product data"
//	@Success		201		{object}	response.Response
//	@Router			/product/create [post]
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

// listProductHandler godoc
//	@Summary		List products
//	@Description	Retrieves a list of products
//	@Tags			product
//	@ID				list-products
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/product/list [get]
func (p productController) listProductHandler(c *gin.Context) {
	resp := p.productService.GetProducts()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// getProductHandler godoc
//	@Summary		Get product
//	@Description	Get single product by id
//	@Tags			product
//	@ID				get-product-by-id
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	response.Response
//	@Router			/product/{id} [get]
func (p productController) getProductHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productService.GetProduct(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// updateProductHandler godoc
//	@Summary		Update product
//	@Description	Update product by id
//	@Tags			product
//	@ID				update-product-by-id
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id		path		int					true	"Product ID"
//	@Param			product	body		utils.UpdateProduct	true	"Update product data"
//	@Success		200		{object}	response.Response
//	@Router			/product/update/{id} [put]
func (p productController) updateProductHandler(c *gin.Context) {
	var product utils.UpdateProduct
	idStr := c.Param("id")

	if err := c.ShouldBindJSON(&product); err != nil {
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
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productService.UpdateProduct(id, product)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// deleteProductHandler godoc
//	@Summary		Delete product
//	@Description	Delete product by id
//	@Tags			product
//	@ID				delete-product-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	response.Response
//	@Router			/product/delete/{id} [delete]
func (p productController) deleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productService.DeleteProduct(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
