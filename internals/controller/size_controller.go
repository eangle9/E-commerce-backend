package controller

import (
	"Eccomerce-website/internals/core/common/router"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/port/service"
	"Eccomerce-website/internals/infra/middleware"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type sizeController struct {
	engine        *router.Router
	sizeService   service.SizeService
	handlerLogger *zap.Logger
}

func NewSizeController(engine *router.Router, sizeService service.SizeService, handlerLogger *zap.Logger) *sizeController {
	return &sizeController{
		engine:        engine,
		sizeService:   sizeService,
		handlerLogger: handlerLogger,
	}
}

func (s *sizeController) InitSizeRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := s.engine
	api := r.Group("/size")

	api.POST("/create", protectedMiddleware(middlewareLogger), s.createSizeHandler)
	api.GET("/list", s.listSizeHandler)
	api.GET("/:id", s.getSizeHandler)
	api.PUT("/update/:id", protectedMiddleware(middlewareLogger), s.updateSizeHandler)
	api.DELETE("/delete/:id", protectedMiddleware(middlewareLogger), s.deleteSizeHandler)
}

// createSizeHandler godoc
// @Summary		     Create size
// @Description	     Insert New product size
// @Tags			 size
// @ID				 create-size
// @Accept			 json
// @Produce		     json
// @Security		 JWT
// @Param			 size	     body		request.SizeRequest	true	"Size data"
// @Success		     201		{object}	response.Response
// @Router			 /size/create [post]
func (s sizeController) createSizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createSizeHandler"),
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
		s.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createSizeHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.SizeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("badRequest",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createSizeHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	if request.Discount.LessThan(request.Price) {
		err := errors.New("discount can't be less than product price")
		errorResponse := entity.AppInternalError.Wrap(err, "discount can't be less than price").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("discount must be greater than or equal to product price",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createSizeHandler"),
			zap.String("requestID", requestID),
			zap.Any("discount", request.Discount),
			zap.Any("price", request.Price),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := s.sizeService.CreateSize(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// listSizeHandler godoc
// @Summary		   List product sizes
// @Description	   Retrieves a list of product sizes
// @Tags		   size
// @ID			   list-product-size
// @Produce		   json
// @Param          page        query   int   false   "Page number"
// @Param          per_page    query   int   false   "Number of items per page"
// @Success		   200	{object}	response.Response
// @Router		   /size/list [get]
func (s sizeController) listSizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listSizeHandler"),
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
		s.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listSizeHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listSizeHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := s.sizeService.GetSizes(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getSizeHandler godoc
// @Summary		  Get size
// @Description	  Get a single size by id
// @Tags		  size
// @ID			  get-size-by-id
// @Produce		  json
// @Param		  id	path		int	true	"Size ID"
// @Success		  200	{object}	response.Response
// @Router		  /size/{id} [get]
func (s sizeController) getSizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSizeHandler"),
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
		s.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSizeHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid size_id.Please enter a valid interger value").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("invalid size_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getSizeHandler"),
			zap.String("requestID", requestID),
			zap.String("sizeID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := s.sizeService.GetSize(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// updateSizeHandler godoc
// @Summary		     Update size
// @Description	     Update product size by id
// @Tags			 size
// @ID				 update-size-by-id
// @Accept			 json
// @Produce		     json
// @Security		 JWT
// @Param			 id		     path		int					true	"Size ID"
// @Param			 size	     body		request.UpdateSize	true	"Update size data"
// @Success		     200		{object}	response.Response
// @Router			 /size/update/{id} [put]
func (s sizeController) updateSizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateSizeHandler"),
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
		s.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateSizeHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid size_id.Please enter a valid interger value").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("invalid size_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateSizeHandler"),
			zap.String("requestID", requestID),
			zap.String("sizeID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var size request.UpdateSize
	if err := c.ShouldBindJSON(&size); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("badRequest",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateSizeHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := s.sizeService.UpdateSize(ctx, id, size, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// deleteSizeHandler godoc
// @Summary		     Delete size
// @Description	     Delete product size by id
// @Tags			 size
// @ID				 delete-size-by-id
// @Produce		     json
// @Security		 JWT
// @Param			 id	path		int	true	"Size ID"
// @Success		     200	        {object}	response.Response
// @Router			 /size/delete/{id} [delete]
func (s sizeController) deleteSizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		s.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteSizeHandler"),
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
		s.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteSizeHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid size_id.Please enter a valid interger value").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		s.handlerLogger.Error("invalid size_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteSizeHandler"),
			zap.String("requestID", requestID),
			zap.String("sizeID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := s.sizeService.DeleteSize(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
