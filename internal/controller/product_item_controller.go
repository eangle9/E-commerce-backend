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
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type productItemController struct {
	engine             *router.Router
	productItemService service.ProductItemService
}

func NewProductItemController(engine *router.Router, productItemService service.ProductItemService) *productItemController {
	return &productItemController{
		engine:             engine,
		productItemService: productItemService,
	}
}

func (p *productItemController) InitProductItemRouter() {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := p.engine
	api := r.Group("/item")

	api.POST("/create", protectedMiddleware(), p.createProductItemHandler)
	api.GET("/list", p.getProductItemsHandler)
	api.GET("/:id", p.getProductItemHandler)
	api.PUT("/update/:id", protectedMiddleware(), p.updateProductItemHandler)
	api.DELETE("/delete/:id", protectedMiddleware(), p.deleteProductItemHandler)
}

// createProductItemHandler godoc
//
//	@Summary		Create product item
//	@Description	insert a new product item
//	@Tags			product item
//	@ID				create-product-item
//	@Accept			mpfd
//	@Produce		json
//	@Security		JWT
//	@Param			product_id		formData	int		true	"Product ID"
//	@Param			color_id		formData	int		false	"Color ID"
//	@Param			price			formData	number	true	"Price"
//	@Param			qty_in_stock	formData	int		false	"Quantity in stock"
//	@Param			image			formData	file	true	"Product Image File"
//	@Success		201				{object}	response.Response
//	@Router			/item/create [post]
func (p productItemController) createProductItemHandler(c *gin.Context) {
	var request request.ProductItemRequest

	productIdStr := c.PostForm("product_id")
	if productIdStr == "" {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "product_id is required",
		}
		c.Set("error", errorResponse)
		return
	}

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid product_item_id.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}
	request.ProductID = productId

	colorIdStr := c.PostForm("color_id")

	if colorIdStr != "" {
		colorId, err := strconv.Atoi(colorIdStr)
		if err != nil {
			errorResponse := response.Response{
				Status:       http.StatusBadRequest,
				ErrorType:    errorcode.InvalidRequest,
				ErrorMessage: "invalid color_id.Please enter a valid integer id",
			}
			c.Set("error", errorResponse)
			return
		}
		request.ColorID = &colorId
	}

	priceStr := c.PostForm("price")
	if priceStr == "" {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "price is required",
		}
		c.Set("error", errorResponse)
		return
	}

	priceInt, err := strconv.Atoi(priceStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid price.Please enter a valid integer id",
		}
		c.Set("error", errorResponse)
		return
	}

	price := decimal.NewFromInt(int64(priceInt))
	request.Price = price

	qtyInStockStr := c.PostForm("qty_in_stock")

	if qtyInStockStr != "" {
		qtyInStock, err := strconv.Atoi(qtyInStockStr)
		if err != nil {
			errorResponse := response.Response{
				Status:       http.StatusBadRequest,
				ErrorType:    errorcode.InvalidRequest,
				ErrorMessage: "invalid qty_in_stock.Please enter a valid integer id",
			}
			c.Set("error", errorResponse)
			return
		}
		request.QtyInStock = &qtyInStock
	}

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

	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.InvalidRequest,
	// 		ErrorMessage: "failed to decode json request body",
	// 	}
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	// validate := c.MustGet("validator").(*validator.Validate)
	// if err := validate.Struct(request); err != nil {
	// 	errorResponse := response.Response{
	// 		Status:       http.StatusBadRequest,
	// 		ErrorType:    errorcode.ValidationError,
	// 		ErrorMessage: customizer7.DecryptErrors(err),
	// 	}
	// 	c.Set("error", errorResponse)
	// 	return
	// }

	resp := p.productItemService.CreateProductItem(request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// getProductItemsHandler godoc
//
//	@Summary		list product items
//	@Description	Retrieves a list of product items
//	@Tags			product item
//	@ID				list-product-item
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/item/list [get]
func (p productItemController) getProductItemsHandler(c *gin.Context) {
	resp := p.productItemService.GetProductItems()
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// getProductItemHandler godoc
//
//	@Summary		Get product item
//	@Description	Get single product item by id
//	@Tags			product item
//	@ID				product-item-by-id
//	@Produce		json
//	@Param			id	path		int	true	"Product item id"
//	@Success		200	{object}	response.Response
//	@Router			/item/{id} [get]
func (p productItemController) getProductItemHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer value",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.GetProductItem(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// updateProductItemHandler godoc
//
//	@Summary		Update product item
//	@Description	Edit product item by id
//	@Tags			product item
//	@ID				update-product-item-by-id
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id		path		int						true	"Product item id"
//	@Param			item	body		utils.UpdateProductItem	true	"Update product item data"
//	@Success		200		{object}	response.Response
//	@Router			/item/update/{id} [put]
func (p productItemController) updateProductItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	var request utils.UpdateProductItem

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
			ErrorMessage: "invalid id.Please enter a valid integer value",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.UpdateProductItem(id, request)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}

// deleteProductItemHandler godoc
//
//	@Summary		Delete product item
//	@Description	Delete product item by id
//	@Tags			product item
//	@ID				delete-product-item-by-id
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		int	true	"Product item id"
//	@Success		200	{object}	response.Response
//	@Router			/item/delete/{id} [delete]
func (p productItemController) deleteProductItemHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: "invalid id.Please enter a valid integer value",
		}
		c.Set("error", errorResponse)
		return
	}

	resp := p.productItemService.DeleteProductItem(id)
	if resp.ErrorType != errorcode.Success {
		c.Set("error", resp)
		return
	}

	c.JSON(resp.Status, resp)
}
