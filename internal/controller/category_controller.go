package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type categoryController struct {
	engine          *router.Router
	categoryService service.ProductCategoryService
	handlerLogger   *zap.Logger
}

func NewCategoryController(engine *router.Router, categoryService service.ProductCategoryService, handlerLogger *zap.Logger) *categoryController {
	return &categoryController{
		engine:          engine,
		categoryService: categoryService,
		handlerLogger:   handlerLogger,
	}
}

func (cat *categoryController) InitCategoryRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := cat.engine
	api := r.Group("/category")

	api.POST("/create", protectedMiddleware(middlewareLogger), cat.createProductCategoryHandler)
	api.GET("/list", cat.listProductCategoryHandler)
	api.GET("/:id", cat.getProductCategoryHandler)
	api.PUT("/update/:id", protectedMiddleware(middlewareLogger), cat.updateProductCategoryHandler)
	api.DELETE("/delete/:id", protectedMiddleware(middlewareLogger), cat.deleteProductCategoryHandler)
}

// createProductCategoryHandler godoc
// @Summary		                Create category
// @Description	                Insert a new product category
// @Tags			            product category
// @ID				            create-product-category
// @Accept			            json
// @Produce		                json
// @Security		            JWT
// @Param			            category	body		request.ProductCategoryRequest	true	"Product category data"
// @Success		                201			{object}	response.Response
// @Router			            /category/create [post]
func (cat categoryController) createProductCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductCategoryHandler"),
			zap.String("context_key", "requestID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	requestID, ok := reqId.(string)
	if !ok {
		err := errors.New("unable to convert type any to string")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductCategoryHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.ProductCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("badRequest",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := cat.categoryService.CreateProductCategory(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// listProductCategoryHandler godoc
// @Summary		              List category
// @Description	              Retrieves a list of product category
// @Tags			          product category
// @ID				          list-product-category
// @Produce		              json
// @Success		              200	{object}	response.Response
// @Router			          /category/list [get]
func (cat categoryController) listProductCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductCategoryHandler"),
			zap.String("context_key", "requestID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	requestID, ok := reqId.(string)
	if !ok {
		err := errors.New("unable to convert type any to string")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductCategoryHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := cat.categoryService.GetProductCategories(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getProductCategoryHandler godoc
// @Summary		             Get category
// @Description	             Get a single product category by id
// @Tags			         product category
// @ID				         get-product-category-by-id
// @Produce		             json
// @Param			         id	path		int	true	"Category ID"
// @Success		             200	{object}	response.Response
// @Router			         /category/{id} [get]
func (cat categoryController) getProductCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductCategoryHandler"),
			zap.String("context_key", "requestID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	requestID, ok := reqId.(string)
	if !ok {
		err := errors.New("unable to convert type any to string")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductCategoryHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("invalid category_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.String("categoryID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := cat.categoryService.GetProductCategory(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// updateProductCategoryHandler godoc
// @Summary		                Update category
// @Description	                Update product category by id
// @Tags			            product category
// @ID				            update-product-category-by-id
// @Accept			            json
// @Produce		                json
// @Security		            JWT
// @Param			            id			path		int						true	"Category ID"
// @Param			            category	body		utils.UpdateCategory	true	"Update product category data"
// @Success		                200			{object}	response.Response
// @Router			            /category/update/{id} [put]
func (cat categoryController) updateProductCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductCategoryHandler"),
			zap.String("context_key", "requestID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	requestID, ok := reqId.(string)
	if !ok {
		err := errors.New("unable to convert type any to string")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductCategoryHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("invalid category_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.String("categoryID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var category utils.UpdateCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("bad request",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := cat.categoryService.UpdateProductCategory(ctx, id, category, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// deleteProductCategoryHandler godoc
// @Summary		                Delete category
// @Description	                Delete product category by id
// @Tags			            product category
// @ID				            delete-product-category-by-id
// @Produce		                json
// @Security		            JWT
// @Param			            id	path		int	true	"Category ID"
// @Success		                200	{object}	response.Response
// @Router			            /category/delete/{id} [delete]
func (cat categoryController) deleteProductCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteProductCategoryHandler"),
			zap.String("context_key", "requestID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	requestID, ok := reqId.(string)
	if !ok {
		err := errors.New("unable to convert type any to string")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		cat.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteProductCategoryHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid id.Please enter a valid integer id").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		cat.handlerLogger.Error("invalid category_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteProductCategoryHandler"),
			zap.String("requestID", requestID),
			zap.String("categoryID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := cat.categoryService.DeleteProductCategory(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
