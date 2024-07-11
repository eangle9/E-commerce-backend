package controller

// import (
// 	"Eccomerce-website/internal/core/common/router"
// 	errorcode "Eccomerce-website/internal/core/entity/error_code"
// 	"Eccomerce-website/internal/core/model/response"
// 	"Eccomerce-website/internal/core/port/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type productsController struct {
// 	engine          *router.Router
// 	productsService service.GetProductService
// }

// func NewProductsController(engine *router.Router, productsService service.GetProductService) *productsController {
// 	return &productsController{
// 		engine:          engine,
// 		productsService: productsService,
// 	}
// }

// func (p *productsController) InitProductsRouter() {
// 	r := p.engine
// 	api := r.Group("/products")

// 	api.GET("/list", p.listProductsHandler)
// 	api.GET("/:id", p.getSingleProductByID)
// }

// // listProductsHandler godoc
// //
// //	@Summary		List of products
// //	@Description	Retrieves a list of products
// //	@Tags			products
// //	@ID				list_of_products
// //	@Produce		json
// //	@Success		200	{object}	response.Response
// //	@Router			/products/list [get]
// func (p productsController) listProductsHandler(c *gin.Context) {
// 	resp := p.productsService.GetAllProducts()
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)
// }

// // getSingleProductByID godoc
// // @Summary		   Get single product
// // @Description	   Get a single product by id
// // @Tags		   products
// // @ID			   get-products-by-id
// // @Produce		   json
// // @Param		   id	path		int	true	"Product ID"
// // @Success		   200	{object}	response.Response
// // @Router		   /products/{id} [get]
// func (p productsController) getSingleProductByID(c *gin.Context) {
// 	idStr := c.Param("id")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		errorResponse := response.Response{
// 			Status:       http.StatusBadRequest,
// 			ErrorType:    errorcode.InvalidRequest,
// 			ErrorMessage: "invalid id.Please enter a valid integer id",
// 		}
// 		c.Set("error", errorResponse)
// 		return
// 	}

// 	resp := p.productsService.GetSingleProduct(id)
// 	if resp.ErrorType != errorcode.Success {
// 		c.Set("error", resp)
// 		return
// 	}

// 	c.JSON(resp.Status, resp)

// }
