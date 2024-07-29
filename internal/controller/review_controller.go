package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/service"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type reviewController struct {
	engine        *router.Router
	reviewService service.ReviewService
	handlerLogger *zap.Logger
}

func NewReviewController(engine *router.Router, reviewService service.ReviewService, handlerLogger *zap.Logger) *reviewController {
	return &reviewController{
		engine:        engine,
		reviewService: reviewService,
		handlerLogger: handlerLogger,
	}
}

func (r *reviewController) InitReviewRouter() {
	router := r.engine
	api := router.Group("/review")

	api.POST("/create", r.createReviewHandler)
	api.GET("/list", r.listReviewHandler)
}

// createReviewHandler godoc
// @Summary		      Create review
// @Description	      Insert a new review
// @Tags			  product review
// @ID				  create-review
// @Accept			  json
// @Produce		      json
// @Security		  JWT
// @Param			  review	body		request.ReviewRequest	true	"review data"
// @Success		      201		{object}	response.Response
// @Router			  /review/create [post]
func (r reviewController) createReviewHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		r.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createReviewHandler"),
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
		r.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createReviewHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.ReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to decode json request body").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		r.handlerLogger.Error("badRequest",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createReviewHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := r.reviewService.CreateReview(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// listReviewHandler godoc
// @Summary		    List review
// @Description	    Retrieves a list of reviews
// @Tags			product review
// @ID				list-review
// @Produce		    json
// @Success		    200	{object}	response.Response
// @Router			/review/list [get]
func (r reviewController) listReviewHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		r.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listReviewHandler"),
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
		r.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listReviewHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		r.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "listReviewHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := r.reviewService.GetReviews(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
