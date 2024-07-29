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

type colorController struct {
	engine        *router.Router
	colorService  service.ColorService
	handlerLogger *zap.Logger
}

func NewColorController(engine *router.Router, colorService service.ColorService, handlerLogger *zap.Logger) *colorController {
	return &colorController{
		engine:        engine,
		colorService:  colorService,
		handlerLogger: handlerLogger,
	}
}

func (color *colorController) InitColorRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := color.engine
	api := r.Group("/color")

	api.POST("/create", protectedMiddleware(middlewareLogger), color.createColorHandler)
	api.GET("/list", protectedMiddleware(middlewareLogger), color.listColorHandler)
	api.GET("/:id", protectedMiddleware(middlewareLogger), color.getColorHandler)
	api.PUT("/update/:id", protectedMiddleware(middlewareLogger), color.updateColorHandler)
	api.DELETE("/delete/:id", protectedMiddleware(middlewareLogger), color.deleteColorHandler)
}

// createColorHandler godoc
// @Summary		      Create color
// @Description	      Insert a new color
// @Tags			  product color
// @ID				  create-color
// @Accept			  json
// @Produce		      json
// @Security		  JWT
// @Param			  color	body		    request.ColorRequest	true	"Color data"
// @Success		      201		{object}	response.Response
// @Router			  /color/create [post]
func (color colorController) createColorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		color.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createColorHandler"),
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
		color.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createColorHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.ColorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		color.handlerLogger.Error("badRequest",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createColorHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := color.colorService.CreateColor(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// listColorHandler godoc
// @Summary		    List color
// @Description	    Retrieves a list of colors
// @Tags			product color
// @ID				list-color
// @Produce		    json
// @Security		JWT
// @Success		    200	{object}	response.Response
// @Router			/color/list [get]
func (color colorController) listColorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		color.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listColorHandler"),
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
		color.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listColorHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		color.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listColorHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := color.colorService.GetColors(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getColorHandler godoc
// @Summary		   Get color
// @Description	   Get a single color by id
// @Tags		   product color
// @ID			   get-color-by-id
// @Produce		   json
// @Security	   JWT
// @Param		   id	path		int	true	"Color ID"
// @Success		   200	{object}	response.Response
// @Router		   /color/{id} [get]
func (color colorController) getColorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		color.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getColorHandler"),
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
		color.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getColorHandler"),
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
		color.handlerLogger.Error("invalid color_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getColorHandler"),
			zap.String("requestID", requestID),
			zap.String("colorID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := color.colorService.GetColor(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// updateColorHandler godoc
// @Summary		      Update color
// @Description	      Update color by id
// @Tags			  product color
// @ID				  update-color-by-id
// @Accept			  json
// @Produce		      json
// @Security		  JWT
// @Param			  id		path		int					true	"Color ID"
// @Param			  color	body		utils.UpdateColor	true	"Update color data"
// @Success		      200		{object}	response.Response
// @Router			  /color/update/{id} [put]
func (color colorController) updateColorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		color.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateColorHandler"),
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
		color.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateColorHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request utils.UpdateColor
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		color.handlerLogger.Error("bad request",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateColorHandler"),
			zap.String("requestID", requestID),
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
		color.handlerLogger.Error("invalid color_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateColorHandler"),
			zap.String("requestID", requestID),
			zap.String("colorID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := color.colorService.UpdateColor(ctx, id, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// deleteColorHandler godoc
// @Summary		      Delete color
// @Description	      Delete color by id
// @Tags			  product color
// @ID				  delete-color-by-id
// @Produce		      json
// @Security		  JWT
// @Param			  id	path		int	true	"Color ID"
// @Success		      200	{object}	response.Response
// @Router			  /color/delete/{id} [delete]
func (color colorController) deleteColorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		color.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteColorHandler"),
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
		color.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteColorHandler"),
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
		return
	}

	resp, err := color.colorService.DeleteColor(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
