package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type sizeController struct {
	engine      *router.Router
	sizeService service.SizeService
}

func NewSizeController(engine *router.Router, sizeService service.SizeService) *sizeController {
	return &sizeController{
		engine:      engine,
		sizeService: sizeService,
	}
}

func (s *sizeController) InitSizeRouter() {
	r := s.engine
	api := r.Group("/size")

	api.POST("/create", s.createSizeHandler)
}

// createSizeHandler godoc
// @Summary          Create size
// @Description      Insert New product size
// @Tags             size
// @ID               create-size
// @Accept           json
// @Produce          json
// @Security         JWT
// @Param            size body request.SizeRequest true "Size data"
// @Success          201 {object} response.Response
// @Router           /size/create [post]
func (s sizeController) createSizeHandler(c *gin.Context) {
	var request request.SizeRequest

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
			ErrorMessage: customizer9.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := s.sizeService.CreateSize(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
