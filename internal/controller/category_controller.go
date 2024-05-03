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
	r := cat.engine
	api := r.Group("/category")

	api.POST("/create", cat.createCategoryHandler)
	api.GET("/list", cat.listCategoryHandler)
	api.GET("/:id", cat.GetProductCategoryHandler)
	api.PUT("/update/:id", cat.updateProductCategoryHandler)
	api.DELETE("/delete/:id", cat.deleteProductCategoryHandler)
}

func (cat categoryController) createCategoryHandler(c *gin.Context) {
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

func (cat categoryController) listCategoryHandler(c *gin.Context) {
	resp := cat.categoryService.GetProductCategories()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (cat categoryController) GetProductCategoryHandler(c *gin.Context) {
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
