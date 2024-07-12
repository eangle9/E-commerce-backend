package main

import (
	_ "Eccomerce-website/docs"
	"Eccomerce-website/internal/controller"
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils/logger"
	"Eccomerce-website/internal/core/common/utils/redis"
	"Eccomerce-website/internal/core/server"
	service "Eccomerce-website/internal/core/service/user_service"

	"Eccomerce-website/internal/infra/config"
	"Eccomerce-website/internal/infra/middleware"
	"Eccomerce-website/internal/infra/repository"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	// "github.com/go-chi/cors"
	"github.com/gin-contrib/cors"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			E-commerce API
//	@version		1.0
//	@description	This is a sample server for an e-commerce platform.

//	@contact.name	Engdawork yismaw
//	@contact.email	engdaworkyismaw9@gmail.com

//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization

//	@host		localhost:9000
//	@schemes	http

func main() {
	appLogger := logger.InitLogger()
	defer appLogger.GetLogger().Sync()

	instance := gin.New()

	// cors middleware
	instance.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://*"}, // Adjust origins as necessary
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}))

	databaseLogger := appLogger.GetLogger().Named("DatabaseLogger")
	serviceLogger := appLogger.GetLogger().Named("ServiceLogger")
	handlerLogger := appLogger.GetLogger().Named("HandlerLogger")
	middlewareLogger := appLogger.GetLogger().Named("MiddlewareLogger")

	errorMiddleware := middleware.ErrorMiddleware
	requestIdMiddleware := middleware.RequestIdMiddleware
	loggerMiddleware := middleware.LoggerMiddleware
	timeoutMiddleware := middleware.TimeoutMiddleware

	redisApp := redis.InitRedis(serviceLogger)
	client := redisApp.GetRedisClient()
	defer client.Close()

	instance.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	instance.Use(gin.Recovery())
	instance.Use(gin.Logger())
	instance.Use(requestIdMiddleware())
	instance.Use(loggerMiddleware(middlewareLogger))
	instance.Use(timeoutMiddleware(middlewareLogger))
	instance.Use(errorMiddleware())

	engine := router.NewRouter(instance)

	dbConfig, httpServerConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDatabase(*dbConfig, databaseLogger)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer db.GetDB().Close()

	// fmt.Println("db: ", db)

	// if err := schema.Migrate(db); err != nil {
	// 	log.Fatal(err)
	// }

	// user service
	userRepo := repository.NewUserRepository(db, databaseLogger)
	userService := service.NewUserService(userRepo, serviceLogger, client)
	userController := controller.NewUserController(engine, userService, handlerLogger)
	userController.InitRouter(middlewareLogger)

	// // product category service
	// categoryRepo := repository.NewProductCategoryRepository(db)
	// categoryService := categoryservice.NewProductCategoryRepository(categoryRepo)
	// categoryController := controller.NewCategoryController(engine, categoryService)
	// categoryController.InitCategoryRouter()

	// // color service
	// colorRepo := repository.NewColorRepository(db)
	// colorService := colorservice.NewColorService(colorRepo)
	// colorController := controller.NewColorController(engine, colorService)
	// colorController.InitColorRouter()

	// // product service
	// productRepo := repository.NewProductRepository(db)
	// productService := productservice.NewProductService(productRepo)
	// productController := controller.NewProductController(engine, productService)
	// productController.InitProductRouter()

	// // product item service
	// productItemRepo := repository.NewProductItemRepository(db)
	// productItemService := productitemservice.NewProductItemService(productItemRepo)
	// productItemController := controller.NewProductItemController(engine, productItemService)
	// productItemController.InitProductItemRouter()

	// // product image service
	// imageRepo := repository.NewProductImageRepository(db)
	// productImageService := productimageservice.NewProductImageService(imageRepo)
	// productImageController := controller.NewProductImageController(engine, productImageService)
	// productImageController.InitProductImageRouter()

	// // cart service
	// cartRepo := repository.NewCartRepository(db)
	// cartService := cartservice.NewCartService(cartRepo)
	// cartController := controller.NewCartController(engine, cartService)
	// cartController.InitCartRouter()

	// // size service
	// sizeRepo := repository.NewSizeRepository(db)
	// sizeService := sizeservice.NewSizeService(sizeRepo)
	// sizeController := controller.NewSizeController(engine, sizeService)
	// sizeController.InitSizeRouter()

	// // products service
	// productsRepo := repository.NewProductsRepository(db)
	// productsService := productsservice.NewProductsService(productsRepo)
	// productsController := controller.NewProductsController(engine, productsService)
	// productsController.InitProductsRouter()

	// // review service
	// reviewRepo := repository.NewReviewRepository(db)
	// reviewService := reviewservice.NewReviewService(reviewRepo)
	// reviewController := controller.NewReviewController(engine, reviewService)
	// reviewController.InitReviewRouter()

	if err := server.Start(instance, *httpServerConfig); err != nil {
		log.Fatal(err)
	}
}
