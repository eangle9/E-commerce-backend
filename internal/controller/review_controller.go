package controller

// import (
// 	"Eccomerce-website/internal/core/common/router"
// 	errorcode "Eccomerce-website/internal/core/entity/error_code"
// 	"Eccomerce-website/internal/core/model/request"
// 	"Eccomerce-website/internal/core/model/response"
// 	"Eccomerce-website/internal/core/port/service"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// )

// type reviewController struct {
// 	engine        *router.Router
// 	reviewService service.ReviewService
// }

// func NewReviewController(engine *router.Router, reviewService service.ReviewService) *reviewController {
// 	return &reviewController{
// 		engine:        engine,
// 		reviewService: reviewService,
// 	}
// }

// func (r *reviewController) InitReviewRouter() {
// 	router := r.engine
// 	api := router.Group("/review")

// 	api.POST("/create", r.createReviewHandler)
// 	api.GET("/list", r.listReviewHandler)
// }

// // createReviewHandler godoc
// // @Summary		      Create review
// // @Description	      Insert a new review
// // @Tags			  product review
// // @ID				  create-review
// // @Accept			  json
// // @Produce		      json
// // @Security		  JWT
// // @Param			  review	body		request.ReviewRequest	true	"review data"
// // @Success		      201		{object}	response.Response
// // @Router			  /review/create [post]
// func (r reviewController) createReviewHandler(c *gin.Context) {
// 	var request request.ReviewRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		response := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: err.Error(),
// 		}
// 		c.Set("error", response)
// 		return
// 	}

// 	validate := c.MustGet("validator").(*validator.Validate)
// 	if err := validate.Struct(request); err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.ValidationError,
// 			ErrorMessage: customizer10.DecryptErrors(err),
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := r.reviewService.CreateReview(request)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // listReviewHandler godoc
// // @Summary		    List review
// // @Description	    Retrieves a list of reviews
// // @Tags			product review
// // @ID				list-review
// // @Produce		    json
// // @Success		    200	{object}	response.Response
// // @Router			/review/list [get]
// func (r reviewController) listReviewHandler(c *gin.Context) {
// 	resp := r.reviewService.GetReviews()
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }
