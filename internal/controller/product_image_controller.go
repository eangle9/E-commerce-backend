package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productImageController struct {
	engine              *router.Router
	productImageService service.ProductImageService
}

func NewProductImageController(engine *router.Router, productImageService service.ProductImageService) *productImageController {
	return &productImageController{
		engine:              engine,
		productImageService: productImageService,
	}
}

func (p *productImageController) InitProductImageRouter() {
	r := p.engine
	api := r.Group("/image")

	api.POST("/upload", p.uploadProductImageHandler)
}

func (p productImageController) uploadProductImageHandler(c *gin.Context) {
	var request request.ProductImageRequest

	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.InvalidRequest,
	// 		ErrorMessage: "failed to decode json request body",
	// 	}
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	productItemIdStr := c.PostForm("product_item_id")
	if productItemIdStr == "" {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "product_item_id is required",
		}
		c.Set("error", errorResponse)
		return
	}

	productItemId, err := strconv.Atoi(productItemIdStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}

	request.ProductItemId = productItemId

	file, err := c.FormFile("image")
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "unable to retrieve file from the upload file",
		}
		c.Set("error", errorResponse)
		return
	}

	if err := c.SaveUploadedFile(file, "./internal/core/common/upload/"+file.Filename); err != nil {
		errorResponse := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: "failed to save image",
		}
		c.Set("error", errorResponse)
		return
	}

	request.File = file

	resp := p.productImageService.CreateProductImage(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
