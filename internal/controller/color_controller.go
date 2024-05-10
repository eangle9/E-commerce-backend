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

type colorController struct {
	engine       *router.Router
	colorService service.ColorService
}

func NewColorController(engine *router.Router, colorService service.ColorService) *colorController {
	return &colorController{
		engine:       engine,
		colorService: colorService,
	}
}

func (color *colorController) InitColorRouter() {
	r := color.engine
	api := r.Group("/color")

	api.POST("/create", color.createColorHandler)
	api.GET("/list", color.ListColorHandler)
	api.GET("/:id", color.GetColorHandler)
	api.PUT("/update/:id", color.updateColorHandler)
	api.DELETE("/delete/:id", color.deleteColorHandler)
}

func (color colorController) createColorHandler(c *gin.Context) {
	var request request.ColorRequest

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
			ErrorMessage: customizer5.DecryptErrors(err),
		}
		c.Set("error", errorResponse)
		return
	}

	resp := color.colorService.CreateColor(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (color colorController) ListColorHandler(c *gin.Context) {
	resp := color.colorService.GetColors()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (color colorController) GetColorHandler(c *gin.Context) {
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

	resp := color.colorService.GetColor(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (color colorController) updateColorHandler(c *gin.Context) {
	var request utils.UpdateColor
	idStr := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "failed to decode json request body",
		}
		c.Set("error", errorResponse)
		return
	}

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

	resp := color.colorService.UpdateColor(id, request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

func (color colorController) deleteColorHandler(c *gin.Context) {
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

	resp := color.colorService.DeleteColor(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
