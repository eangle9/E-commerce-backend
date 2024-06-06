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

type categoryController struct {
	engine          *router.Router
	categoryService service.ProductCategoryService
}

func NewCategoryController(engine *router.Router, categoryService service.ProductCategoryService) *categoryController {
	return &categoryController{
		engine:          engine,
		categoryService: categoryService,
	}
}

func (cat *categoryController) InitCategoryRouter() {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := cat.engine
	api := r.Group("/category")

	api.POST("/create", protectedMiddleware(), cat.createProductCategoryHandler)
	api.GET("/list", cat.listProductCategoryHandler)
	api.GET("/:id", cat.getProductCategoryHandler)
	api.PUT("/update/:id", protectedMiddleware(), cat.updateProductCategoryHandler)
	api.DELETE("/delete/:id", protectedMiddleware(), cat.deleteProductCategoryHandler)
}

// createProductCategoryHandler godoc
//	@Summary		Create category
//	@Description	Insert a new product category
//	@Tags			product category
//	@ID				create-product-category
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			category	body		request.ProductCategoryRequest	true	"Product category data"
//	@Success		201			{object}	response.Response
//	@Router			/category/create [post]
func (cat categoryController) createProductCategoryHandler(c *gin.Context) {
	var request request.ProductCategoryRequest

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
			ErrorMessage: customizer4.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := cat.categoryService.CreateProductCategory(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// listProductCategoryHandler godoc
//	@Summary		List category
//	@Description	Retrieves a list of product category
//	@Tags			product category
//	@ID				list-product-category
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/category/list [get]
func (cat categoryController) listProductCategoryHandler(c *gin.Context) {
	resp := cat.categoryService.GetProductCategories()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// getProductCategoryHandler godoc
//	@Summary		Get category
//	@Description	Get a single product category by id
//	@Tags			product category
//	@ID				get-product-category-by-id
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	response.Response
//	@Router			/category/{id} [get]
func (cat categoryController) getProductCategoryHandler(c *gin.Context) {
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

	resp := cat.categoryService.GetProductCategory(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// updateProductCategoryHandler godoc
//	@Summary		Update category
//	@Description	Update product category by id
//	@Tags			product category
//	@ID				update-product-category-by-id
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id			path		int						true	"Category ID"
//	@Param			category	body		utils.UpdateCategory	true	"Update product category data"
//	@Success		200			{object}	response.Response
//	@Router			/category/update/{id} [put]
func (cat categoryController) updateProductCategoryHandler(c *gin.Context) {
	idStr := c.Param("id")
	var category utils.UpdateCategory

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

	if err := c.ShouldBindJSON(&category); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := cat.categoryService.UpdateProductCategory(id, category)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// deleteProductCategoryHandler godoc
//	@Summary		Delete category
//	@Description	Delete product category by id
//	@Tags			product category
//	@ID				delete-product-category-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	response.Response
//	@Router			/category/delete/{id} [delete]
func (cat categoryController) deleteProductCategoryHandler(c *gin.Context) {
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

	resp := cat.categoryService.DeleteProductCategory(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
