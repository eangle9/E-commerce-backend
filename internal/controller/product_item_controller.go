package controller

import (
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils"
	validationdata "Eccomerce-website/internal/core/common/utils/validationData"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/infra/middleware"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type productItemController struct {
	engine             *router.Router
	productItemService service.ProductItemService
	handlerLogger      *zap.Logger
}

func NewProductItemController(engine *router.Router, productItemService service.ProductItemService, handlerLogger *zap.Logger) *productItemController {
	return &productItemController{
		engine:             engine,
		productItemService: productItemService,
		handlerLogger:      handlerLogger,
	}
}

func (p *productItemController) InitProductItemRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := p.engine
	api := r.Group("/item")

	api.POST("/create", protectedMiddleware(middlewareLogger), p.createProductItemHandler)
	api.GET("/list", p.getProductItemsHandler)
	api.GET("/:id", p.getProductItemHandler)
	api.PUT("/update/:id", protectedMiddleware(middlewareLogger), p.updateProductItemHandler)
	api.DELETE("/delete/:id", protectedMiddleware(middlewareLogger), p.deleteProductItemHandler)
}

// createProductItemHandler godoc
// @Summary		            Create product item
// @Description	            insert a new product item
// @Tags			        product item
// @ID				        create-product-item
// @Accept			        mpfd
// @Produce		            json
// @Security		        JWT
// @Param                   product_id      formData    int     true    "Product ID"
// @Param			        color_id		formData	int		false	"Color ID"
// @Param			        price			formData	number	true	"Price"
// @Param                   discount        formData    number  false   "Discount"
// @Param			        qty_in_stock	formData	int		false	"Quantity in stock"
// @Param			        file			formData	file	true	"Product Image File"
// @Success		            201				{object}	response.Response
// @Router			        /item/create [post]
func (p productItemController) createProductItemHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
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
			zap.String("function", "createProductItemHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request request.ProductItemRequest
	productIdStr := c.PostForm("product_id")
	if productIdStr == "" {
		err := errors.New("product_id can't be blank")
		errorResponse := entity.BadRequest.Wrap(err, "product_id is required").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("product_id is required",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "invalid product_id.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid product_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("productID", productIdStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}
	request.ProductID = productId

	colorIdStr := c.PostForm("color_id")
	if colorIdStr != "" {
		colorId, err := strconv.Atoi(colorIdStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid color_id.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("invalid color_id",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "createProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("colorID", colorIdStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}
		request.ColorID = &colorId
	}

	priceStr := c.PostForm("price")
	if priceStr == "" {
		err := errors.New("price can't be blank")
		errorResponse := entity.BadRequest.Wrap(err, "price is required").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("price is required",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "error converting string to decimal").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid price",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("price", priceStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}
	request.Price = price

	discountStr := c.PostForm("discount")
	if discountStr != "" {
		discount, err := decimal.NewFromString(discountStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "error converting string to decimal").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("invalid discount",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "createProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("discount", discountStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}
		request.Discount = discount
	}

	qtyInStockStr := c.PostForm("qty_in_stock")
	if qtyInStockStr == "" {
		err := errors.New("qtyInStock can't be blank")
		errorResponse := entity.BadRequest.Wrap(err, "qtyInStock is required").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("qtyInStock is required",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	qtyInStock, err := strconv.Atoi(qtyInStockStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "invalid qty_in_stock.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid qtyInStock",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("qtyInStock", qtyInStockStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}
	request.QtyInStock = qtyInStock

	file, err := c.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			errorResponse := entity.BadRequest.Wrap(err, "no file was uploaded").WithProperty(entity.StatusCode, 400)
			c.Error(errorResponse)
			p.handlerLogger.Error("upload file can't be blank",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "createProductItemHandler"),
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
				zap.String("function", "createProductItemHandler"),
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
		p.handlerLogger.Error("unable to save the upload file",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "createProductItemHandler"),
			zap.String("requestID", requestID),
			zap.Any("file", file),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	request.File = file

	resp, err := p.productItemService.CreateProductItem(ctx, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getProductItemsHandler godoc
// @Summary		          list product items
// @Description	          Retrieves a list of product items
// @Tags			      product item
// @ID				      list-product-item
// @Produce		          json
// @Success		          200	{object}	response.Response
// @Router			      /item/list [get]
func (p productItemController) getProductItemsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductItemsHandler"),
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
			zap.String("function", "getProductItemsHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var paginationQuery request.PaginationQuery
	if err := c.ShouldBindQuery(&paginationQuery); err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "failed to bind pagination query to struct").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("failed to bind pagination query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductItemsHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := p.productItemService.GetProductItems(ctx, paginationQuery, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// getProductItemHandler godoc
// @Summary		         Get product item
// @Description	         Get single product item by id
// @Tags			     product item
// @ID				     product-item-by-id
// @Produce		         json
// @Param			     id	path		int	true	"Product item id"
// @Success		         200	        {object}	response.Response
// @Router		      	 /item/{id} [get]
func (p productItemController) getProductItemHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductItemHandler"),
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
			zap.String("function", "getProductItemHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.BadRequest.Wrap(err, "invalid product_item_id.Please enter a valid interger value").WithProperty(entity.StatusCode, 400)
		c.Error(errorResponse)
		p.handlerLogger.Error("invalid product_item_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "getProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("productItemID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := p.productItemService.GetProductItem(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// updateProductItemHandler godoc
// @Summary		            Update product item
// @Description	            Edit product item by id
// @Tags			        product item
// @ID				        update-product-item-by-id
// @Accept			        json
// @Produce		            json
// @Security		        JWT
// @Param			        id		        path		int		true	"Product item id"
// @Param			        product_id		formData	int		true	"Product ID"
// @Success		            200		{object}	response.Response
// @Router			        /item/update/{id} [put]
func (p productItemController) updateProductItemHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductItemHandler"),
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
			zap.String("function", "updateProductItemHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	var request utils.UpdateProductItem
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "invalid id.Please enter a valid integer value").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("incorrect product_item_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("productItemID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	productIdStr := c.PostForm("product_id")
	if productIdStr != "" {
		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid product_id.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("incorrect product_id",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("productID", productIdStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}
		request.ProductID = &productId
	}

	colorIdStr := c.PostForm("color_id")
	if colorIdStr != "" {
		colorId, err := strconv.Atoi(colorIdStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid color_id.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("incorrect color_id",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("colorID", colorIdStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}
		request.ColorID = &colorId
	}

	priceStr := c.PostForm("price")
	if priceStr != "" {
		priceInt, err := strconv.Atoi(priceStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid price.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("incorrect price",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("price", priceStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		price := decimal.NewFromInt(int64(priceInt))
		request.Price = price
	}

	discountStr := c.PostForm("discount")
	if discountStr != "" {
		discountInt, err := strconv.Atoi(discountStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid discount.Please enter a valid integer value").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("incorrect discount",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("discount", discountStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		discount := decimal.NewFromInt(int64(discountInt))
		request.Discount = discount
	}

	qtyInStockStr := c.PostForm("qty_in_stock")
	if qtyInStockStr != "" {
		qtyInStock, err := strconv.Atoi(qtyInStockStr)
		if err != nil {
			errorResponse := entity.AppInternalError.Wrap(err, "invalid qty_in_stock.Please enter a valid integer id").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("incorrect qtyInStcok",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.String("qtyInStcok", qtyInStockStr),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}
		request.QtyInStock = &qtyInStock
	}

	file, err := c.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		errorResponse := entity.AppInternalError.Wrap(err, "unable to retrieve file from the upload file").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("failed to retrieve file from the upload file",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "updateProductItemHandler"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	if err == nil {
		if err := validationdata.ImageFileValidation(file, p.handlerLogger, requestID); err != nil {
			c.Error(err)
			return
		}

		if err := c.SaveUploadedFile(file, "./internal/core/common/upload/"+file.Filename); err != nil {
			errorResponse := entity.UnableToSaveFile.Wrap(err, "failed to save upload file").WithProperty(entity.StatusCode, 500)
			c.Error(errorResponse)
			p.handlerLogger.Error("unable to save image file",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "handlerLayer"),
				zap.String("function", "updateProductItemHandler"),
				zap.String("requestID", requestID),
				zap.Any("file", file),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		request.File = file
	}

	resp, err := p.productItemService.UpdateProductItem(ctx, id, request, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}

// deleteProductItemHandler godoc
// @Summary		            Delete product item
// @Description	            Delete product item by id
// @Tags			        product item
// @ID				        delete-product-item-by-id
// @Produce		            json
// @Security		        JWT
// @Param			        id	   path		   int	true	"Product item id"
// @Success		            200	   {object}	   response.Response
// @Router			        /item/delete/{id} [delete]
func (p productItemController) deleteProductItemHandler(c *gin.Context) {
	ctx := c.Request.Context()

	reqId, exist := c.Get("requestID")
	if !exist {
		err := errors.New("unable to get requestID from the gin context")
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get request id").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("requestID is not exist in the gin context",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteProductItemHandler"),
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
			zap.String("function", "deleteProductItemHandler"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "invalid id.Please enter a valid integer value").WithProperty(entity.StatusCode, 500)
		c.Error(errorResponse)
		p.handlerLogger.Error("incorrect product_item_id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "deleteProductItemHandler"),
			zap.String("requestID", requestID),
			zap.String("productItemID", idStr),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return
	}

	resp, err := p.productItemService.DeleteProductItem(ctx, id, requestID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(resp.StatusCode, resp)
}
