package controller

// import (
// 	"Eccomerce-website/internal/core/common/router"
// 	"Eccomerce-website/internal/core/common/utils"
// 	errorcode "Eccomerce-website/internal/core/entity/error_code"
// 	"Eccomerce-website/internal/core/model/request"
// 	"Eccomerce-website/internal/core/model/response"
// 	"Eccomerce-website/internal/core/port/service"
// 	"Eccomerce-website/internal/infra/middleware"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// )

// type sizeController struct {
// 	engine      *router.Router
// 	sizeService service.SizeService
// }

// func NewSizeController(engine *router.Router, sizeService service.SizeService) *sizeController {
// 	return &sizeController{
// 		engine:      engine,
// 		sizeService: sizeService,
// 	}
// }

// func (s *sizeController) InitSizeRouter() {
// 	protectedMiddleware := middleware.ProtectedMiddleware
// 	r := s.engine
// 	api := r.Group("/size")

// 	api.POST("/create", protectedMiddleware(), s.createSizeHandler)
// 	api.GET("/list", s.listSizeHandler)
// 	api.GET("/:id", s.getSizeHandler)
// 	api.PUT("/update/:id", protectedMiddleware(), s.updateSizeHandler)
// 	api.DELETE("/delete/:id", protectedMiddleware(), s.deleteSizeHandler)
// }

// // createSizeHandler godoc
// //
// //	@Summary		Create size
// //	@Description	Insert New product size
// //	@Tags			size
// //	@ID				create-size
// //	@Accept			json
// //	@Produce		json
// //	@Security		JWT
// //	@Param			size	body		request.SizeRequest	true	"Size data"
// //	@Success		201		{object}	response.Response
// //	@Router			/size/create [post]
// func (s sizeController) createSizeHandler(c *gin.Context) {
// 	var request request.SizeRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "failed to decode json request body",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	if request.Discount.LessThan(request.Price) {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "discount can't be less than product price",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	validate := c.MustGet("validator").(*validator.Validate)
// 	if err := validate.Struct(request); err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.ValidationError,
// 			ErrorMessage: customizer9.DecryptErrors(err),
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := s.sizeService.CreateSize(request)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // listSizeHandler godoc
// //
// //	@Summary		List product sizes
// //	@Description	Retrieves a list of product sizes
// //	@Tags			size
// //	@ID				list-product-size
// //	@Produce		json
// //	@Success		200	{object}	response.Response
// //	@Router			/size/list [get]
// func (s sizeController) listSizeHandler(c *gin.Context) {
// 	resp := s.sizeService.GetSizes()
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // getSizeHandler godoc
// //
// //	@Summary		Get size
// //	@Description	Get a single size by id
// //	@Tags			size
// //	@ID				get-size-by-id
// //	@Produce		json
// //	@Param			id	path		int	true	"Size ID"
// //	@Success		200	{object}	response.Response
// //	@Router			/size/{id} [get]
// func (s sizeController) getSizeHandler(c *gin.Context) {
// 	idStr := c.Param("id")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "invalid size_id.Please enter a valid interger value",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := s.sizeService.GetSize(id)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // updateSizeHandler godoc
// //
// //	@Summary		Update size
// //	@Description	Update product size by id
// //	@Tags			size
// //	@ID				update-size-by-id
// //	@Accept			json
// //	@Produce		json
// //	@Security		JWT
// //	@Param			id		path		int					true	"Size ID"
// //	@Param			size	body		utils.UpdateSize	true	"Update size data"
// //	@Success		200		{object}	response.Response
// //	@Router			/size/update/{id} [put]
// func (s sizeController) updateSizeHandler(c *gin.Context) {
// 	var size utils.UpdateSize
// 	idStr := c.Param("id")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "invalid size_id.Please enter a valid interger value",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&size); err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "failed to decode json request body",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := s.sizeService.UpdateSize(id, size)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // deleteSizeHandler godoc
// //
// //	@Summary		Delete size
// //	@Description	Delete product size by id
// //	@Tags			size
// //	@ID				delete-size-by-id
// //	@Produce		json
// //	@Security		JWT
// //	@Param			id	path		int	true	"Size ID"
// //	@Success		200	{object}	response.Response
// //	@Router			/size/delete/{id} [delete]
// func (s sizeController) deleteSizeHandler(c *gin.Context) {
// 	idStr := c.Param("id")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "invalid size_id.Please enter a valid interger value",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := s.sizeService.DeleteSize(id)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }
