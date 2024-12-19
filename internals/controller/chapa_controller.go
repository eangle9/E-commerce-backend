package controller

import (
	"Eccomerce-website/internals/core/common/router"
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/infra/middleware"
	platform "Eccomerce-website/platform_appp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type chapaController struct {
	engine        *router.Router
	chapaService  platform.API
	handlerLogger *zap.Logger
}

func NewChapaController(engine *router.Router, chapaService platform.API, handlerLogger *zap.Logger) *chapaController {
	return &chapaController{
		engine:        engine,
		chapaService:  chapaService,
		handlerLogger: handlerLogger,
	}
}

func (chapa *chapaController) InitChapaRouter(middlewareLogger *zap.Logger) {
	protectedMiddleware := middleware.ProtectedMiddleware
	r := chapa.engine
	api := r.Group("/chapa")

	api.POST("/init", protectedMiddleware(middlewareLogger), chapa.initPaymentHandler)
	api.GET("/verify", protectedMiddleware(middlewareLogger), chapa.verifyPayment)
}

func (chapa *chapaController) initPaymentHandler(c *gin.Context) {
	var paymentRequest dto.PaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		err = entity.BadRequest.Wrap(err, "failed to bind request body").WithProperty(entity.StatusCode, 400)
		c.Error(err)
		return
	}

	res, err := chapa.chapaService.InitiatePayment(&paymentRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}

func (chapa *chapaController) verifyPayment(c *gin.Context) {
	var txRef string
	if err := c.ShouldBindJSON(&txRef); err != nil {
		err := entity.BadRequest.Wrap(err, "unable to bind json request body").WithProperty(entity.StatusCode, 400)
		c.Error(err)
		return
	}
	res, err := chapa.chapaService.VerifyPayment(txRef)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, res)
}
