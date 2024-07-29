package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/service"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type productsController struct {
	engine          *router.Router
	productsService service.GetProductService
	handlerLogger   *zap.Logger
}

func NewProductsController(engine *router.Router, productsService service.GetProductService, handlerLogger *zap.Logger) *productsController {
	return &productsController{
		engine:          engine,
		productsService: productsService,
		handlerLogger:   handlerLogger,
	}
}

func (p *productsController) InitProductsRouter() {
	r := p.engine
	api := r.Group("/products")

	api.GET("/list", p.listProductsHandler)
	api.GET("/:id", p.getSingleProductByID)
}

// listProductsHandler godoc
// @Summary		       List of products
// @Description	       Retrieves a list of products
// @Tags			   products
// @ID				   list_of_products
// @Produce		       json
// @Success		       200	{object}	response.Response
// @Router			   /products/list [get]
func (p productsController) listProductsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductsHandler"),
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
		p.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductsHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listProductsHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := p.productsService.GetAllProducts(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getSingleProductByID godoc
// @Summary		        Get single product
// @Description	        Get a single product by id
// @Tags		        products
// @ID			        get-products-by-id
// @Produce		        json
// @Param		        id	path		int	true	"Product ID"
// @Success		        200	{object}	response.Response
// @Router		        /products/{id} [get]
func (p productsController) getSingleProductByID(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSingleProductByID"),
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
		p.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSingleProductByID"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid product_id.Please enter a valid interger value").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid product_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSingleProductByID"),
			zap.String("requestID", requestID),
			zap.String("productID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := p.productsService.GetSingleProduct(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)

}
