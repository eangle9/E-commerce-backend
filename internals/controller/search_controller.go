package controller

// import (
// 	"Eccomerce-website/internal/core/common/router"
// 	"Eccomerce-website/internal/core/entity"
// 	"Eccomerce-website/internal/core/model/request"
// 	"Eccomerce-website/internal/core/port/service"
// 	"errors"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// type searchController struct {
// 	engine          *router.Router
// 	productsService service.SearchService
// 	handlerLogger   *zap.Logger
// }

// func NewSearchController(engine *router.Router, searchService service.SearchService, handlerLogger *zap.Logger) *searchController {
// 	return &searchController{
// 		engine:          engine,
// 		productsService: searchService,
// 		handlerLogger:   handlerLogger,
// 	}
// }

// func (p *searchController) InitSearchRouter() {
// 	r := p.engine
// 	api := r.Group("/products")

// 	api.GET("/list", p.listSearchHandler)
// }

// // listProductsHandler godoc
// // @Summary		       List of products
// // @Description	       Retrieves a list of products
// // @Tags			   products
// // @ID				   list_of_products
// // @Produce		       json
// // @Param              page       query   int   false   "Page number"
// // @Param              per_page   query   int   false   "Number of items per page"
// // @Success		       200	{object}	response.Response
// // @Router			   /products/list [get]
// func (p searchController) listSearchHandler(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	reqId, exist := c.Get("requestID")
// 	if !exist {
// 		err := errors.New("unable to get requestID from the gin context")
// 		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("requestID is not exist in the gin context",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.String("context_key", "requestID"),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	requestID, ok := reqId.(string)
// 	if !ok {
// 		err := errors.New("unable to convert type any to string")
// 		errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string").WithProperty(entity.StatusCode, 500)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("requestId is not exist in the context",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	var paginationQuery request.PaginationQuery
// 	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
// 		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("failed to bind pagination query",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.String("requestID", requestID),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	var searchQuery request.SearchQuery
// 	if err := c.ShouldBindQuery(&searchQuery); err != nil {
// 		errorResponse := entity.BadRequest.Wrap(err, "failed to bind search query to struct").WithProperty(entity.StatusCode, 400)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("failed to bind search query",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.String("requestID", requestID),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	var categoryQuery request.CategoryQuery
// 	if err := c.ShouldBindQuery(&categoryQuery); err != nil {
// 		errorResponse := entity.BadRequest.Wrap(err, "failed to bind category query to struct").WithProperty(entity.StatusCode, 400)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("failed to bind category query",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.String("requestID", requestID),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	var sortQuery request.SortQuery
// 	if err := c.ShouldBindQuery(&sortQuery); err != nil {
// 		errorResponse := entity.BadRequest.Wrap(err, "failed to bind sort query to struct").WithProperty(entity.StatusCode, 400)
// 		c.Error(errorResponse)
// 		p.handlerLogger.Error("failed to bind sort query",
// 			zap.String("timestamp", time.Now().Format(time.RFC3339)),
// 			zap.String("layer", "handlerLayer"),
// 			zap.String("function", "listProductsHandler"),
// 			zap.String("requestID", requestID),
// 			zap.Error(errorResponse),
// 			zap.Stack("stacktrace"),
// 		)
// 		return
// 	}

// 	resp, err := p.productsService.GetAllProducts(ctx, paginationQuery, searchQuery, categoryQuery, sortQuery, requestID)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(resp.StatusCode, resp)
// }
