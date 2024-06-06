package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
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
	protectedMiddleware := middleware.ProtectedMiddleware
	r := color.engine
	api := r.Group("/color")

	api.POST("/create", protectedMiddleware(), color.createColorHandler)
	api.GET("/list", protectedMiddleware(), color.listColorHandler)
	api.GET("/:id", protectedMiddleware(), color.getColorHandler)
	api.PUT("/update/:id", protectedMiddleware(), color.updateColorHandler)
	api.DELETE("/delete/:id", protectedMiddleware(), color.deleteColorHandler)
}

// createColorHandler godoc
//	@Summary		Create color
//	@Description	Insert a new color
//	@Tags			product color
//	@ID				create-color
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			color	body		request.ColorRequest	true	"Color data"
//	@Success		201		{object}	response.Response
//	@Router			/color/create [post]
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

// listColorHandler godoc
//	@Summary		List color
//	@Description	Retrieves a list of colors
//	@Tags			product color
//	@ID				list-color
//	@Produce		json
//	@Security		JWT
//	@Success		200	{object}	response.Response
//	@Router			/color/list [get]
func (color colorController) listColorHandler(c *gin.Context) {
	resp := color.colorService.GetColors()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// getColorHandler godoc
//	@Summary		Get color
//	@Description	Get a single color by id
//	@Tags			product color
//	@ID				get-color-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Color ID"
//	@Success		200	{object}	response.Response
//	@Router			/color/{id} [get]
func (color colorController) getColorHandler(c *gin.Context) {
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

// updateColorHandler godoc
//	@Summary		Update color
//	@Description	Update color by id
//	@Tags			product color
//	@ID				update-color-by-id
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id		path		int					true	"Color ID"
//	@Param			color	body		utils.UpdateColor	true	"Update color data"
//	@Success		200		{object}	response.Response
//	@Router			/color/update/{id} [put]
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

// deleteColorHandler godoc
//	@Summary		Delete color
//	@Description	Delete color by id
//	@Tags			product color
//	@ID				delete-color-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Color ID"
//	@Success		200	{object}	response.Response
//	@Router			/color/delete/{id} [delete]
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
