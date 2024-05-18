package controller

import (
	"Eccomerce-website/internal/core/common/router"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"net/http"
	"path/filepath"
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
	protectedMiddleware := middleware.ProtectedMiddleware
	r := p.engine
	api := r.Group("/image")

	api.POST("/upload", protectedMiddleware(), p.uploadProductImageHandler)
}

// uploadProductImageHandler godoc
//
//	@Summary		Upload image
//	@Description	Upload an image for a product in cloudinary
//	@Tags			product image
//	@ID				upload-product-image
//	@Accept			mpfd
//	@Produce		json
//	@Security		JWT
//	@Param			product_item_id	formData	int		true	"Product Item ID"
//	@Param			image			formData	file	true	"product image"
//	@Success		200				{object}	response.Response
//	@Router			/image/upload [post]
func (p productImageController) uploadProductImageHandler(c *gin.Context) {
	var request request.ProductImageRequest

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
			ErrorMessage: "invalid product_item_id.Please enter a valid integer id",
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

	maxUploadSize := 8 * 1024 * 1024
	fileSize := file.Size

	if fileSize > int64(maxUploadSize) {
		errorResponse := response.Response{
			Status:       http.StatusRequestEntityTooLarge,
			ErrorType:    "FILE_TOO_LARGE",
			ErrorMessage: "the uploaded product image is too large.Please upload a size less than 8MB",
		}
		c.Set("error", errorResponse)
		return
	}

	validExt := map[string]bool{
		".jpeg": true,
		".png":  true,
		".jpg":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !validExt[ext] {
		errorResponse := response.Response{
			Status:       http.StatusUnsupportedMediaType,
			ErrorType:    "INVALID_FILE_EXTENSION",
			ErrorMessage: "invalid file extension",
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
