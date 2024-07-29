package main

import (
	_ "Eccomerce-website/docs"
	"Eccomerce-website/internal/controller"
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/common/utils/logger"
	"Eccomerce-website/internal/core/common/utils/redis"
	"Eccomerce-website/internal/core/server"
	categoryservice "Eccomerce-website/internal/core/service/category_service"
	colorservice "Eccomerce-website/internal/core/service/color_service"
	productimageservice "Eccomerce-website/internal/core/service/product_image_service"
	productitemservice "Eccomerce-website/internal/core/service/product_item_service"
	productservice "Eccomerce-website/internal/core/service/product_service"
	productsservice "Eccomerce-website/internal/core/service/products_service"
	reviewservice "Eccomerce-website/internal/core/service/review_service"
	sizeservice "Eccomerce-website/internal/core/service/size_service"
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

	// product category service
	categoryRepo := repository.NewProductCategoryRepository(db, databaseLogger)
	categoryService := categoryservice.NewProductCategoryRepository(categoryRepo, serviceLogger)
	categoryController := controller.NewCategoryController(engine, categoryService, handlerLogger)
	categoryController.InitCategoryRouter(middlewareLogger)

	// color service
	colorRepo := repository.NewColorRepository(db, databaseLogger)
	colorService := colorservice.NewColorService(colorRepo, serviceLogger)
	colorController := controller.NewColorController(engine, colorService, handlerLogger)
	colorController.InitColorRouter(middlewareLogger)

	// product service
	productRepo := repository.NewProductRepository(db, databaseLogger)
	productService := productservice.NewProductService(productRepo, serviceLogger)
	productController := controller.NewProductController(engine, productService, handlerLogger)
	productController.InitProductRouter(middlewareLogger)

	// product item service
	productItemRepo := repository.NewProductItemRepository(db, databaseLogger)
	productItemService := productitemservice.NewProductItemService(productItemRepo, serviceLogger)
	productItemController := controller.NewProductItemController(engine, productItemService, handlerLogger)
	productItemController.InitProductItemRouter(middlewareLogger)

	// product image service
	imageRepo := repository.NewProductImageRepository(db, databaseLogger)
	productImageService := productimageservice.NewProductImageService(imageRepo, serviceLogger)
	productImageController := controller.NewProductImageController(engine, productImageService, handlerLogger)
	productImageController.InitProductImageRouter(middlewareLogger)

	// // cart service
	// cartRepo := repository.NewCartRepository(db)
	// cartService := cartservice.NewCartService(cartRepo)
	// cartController := controller.NewCartController(engine, cartService)
	// cartController.InitCartRouter()

	// size service
	sizeRepo := repository.NewSizeRepository(db, databaseLogger)
	sizeService := sizeservice.NewSizeService(sizeRepo, serviceLogger)
	sizeController := controller.NewSizeController(engine, sizeService, handlerLogger)
	sizeController.InitSizeRouter(middlewareLogger)

	// products service
	productsRepo := repository.NewProductsRepository(db, databaseLogger)
	productsService := productsservice.NewProductsService(productsRepo, serviceLogger)
	productsController := controller.NewProductsController(engine, productsService, handlerLogger)
	productsController.InitProductsRouter()

	// review service
	reviewRepo := repository.NewReviewRepository(db, databaseLogger)
	reviewService := reviewservice.NewReviewService(reviewRepo, serviceLogger)
	reviewController := controller.NewReviewController(engine, reviewService, handlerLogger)
	reviewController.InitReviewRouter()

	if err := server.Start(instance, *httpServerConfig); err != nil {
		log.Fatal(err)
	}
}
