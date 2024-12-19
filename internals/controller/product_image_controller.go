package controller

import (
	"Eccomerce-website/internals/core/common/router"
	validationdata "Eccomerce-website/internals/core/common/utils/validationData"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/port/service"
	"Eccomerce-website/internals/infra/middleware"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type productImageController struct {
	engine              *router.Router
	productImageService service.ProductImageService
	handlerLogger       *zap.Logger
}

func NewProductImageController(engine *router.Router, productImageService service.ProductImageService, handlerLogger *zap.Logger) *productImageController {
	return &productImageController{
		engine:              engine,
		productImageService: productImageService,
		handlerLogger:       handlerLogger,
	}
}

func (p *productImageController) InitProductImageRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := p.engine
	api := r.Group("/image")

	api.POST("/upload", protectedMiddleware(middlewareLogger), p.uploadProductImageHandler)
}

// uploadProductImageHandler godoc
// @Summary		             Upload image
// @Description	             Upload an image for a product in cloudinary
// @Tags			         product image
// @ID				         upload-product-image
// @Accept			         mpfd
// @Produce		             json
// @Security		         JWT
// @Param			         product_item_id	formData	int		true	"Product Item ID"
// @Param			         image			    formData	file	true	"product image"
// @Success		             200				{object}	response.Response
// @Router			         /image/upload [post]
func (p productImageController) uploadProductImageHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "uploadProductImageHandler"),
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
		p.handlerLogger.Error("requestId is not exist in the context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "uploadProductImageHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.ProductImageRequest
	productItemIdStr := c.PostForm("product_item_id")
	if productItemIdStr == "" {
		err := errors.New("product_item_id can't be blank")
		errorResponse := entity.BadRequest.Wrap(err, "product_item_id is required").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("product_item_id is required",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "uploadProductImageHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	productItemId, err := strconv.Atoi(productItemIdStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "invalid product_item_id.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid product_item_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "uploadProductImageHandler"),
			zap.String("requestID", requestID),
			zap.String("productItemID", productItemIdStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	request.ProductItemId = productItemId

	file, err := c.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			errorResponse := entity.BadRequest.Wrap(err, "no file was uploaded").WithProperty(entity.StatusCode, 400)
			c.Error(errorResponse)
			p.handlerLogger.Error("upload file is required",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "uploadProductImageHandler"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		} else {
			errorResponse := entity.AppInternalError.Wrap(err, "unable to retrieve file from the upload file").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("failed to retrieve file from the upload file",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "uploadProductImageHandler"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

	}

	if err := validationdata.ImageFileValidation(file, p.handlerLogger, requestID); err != nil {
		c.Error(err)
		return
	}

	if err := c.SaveUploadedFile(file, "./internal/core/common/upload/"+file.Filename); err != nil {
		errorResponse := entity.UnableToSaveFile.Wrap(err, "failed to save image file").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("unable to save upload image",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "uploadProductImageHandler"),
			zap.String("requestID", requestID),
			zap.Any("file", file),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	request.File = file

	resp, err := p.productImageService.CreateProductImage(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
