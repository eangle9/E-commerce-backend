package main

import (
	_ "Eccomerce-website/docs"
	"Eccomerce-website/internal/controller"
	"Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/core/server"
	cartservice "Eccomerce-website/internal/core/service/cart_service"
	categoryservice "Eccomerce-website/internal/core/service/category_service"
	colorservice "Eccomerce-website/internal/core/service/color_service"
	productimageservice "Eccomerce-website/internal/core/service/product_image_service"
	productitemservice "Eccomerce-website/internal/core/service/product_item_service"
	productservice "Eccomerce-website/internal/core/service/product_service"
	productsservice "Eccomerce-website/internal/core/service/products_service"
	sizeservice "Eccomerce-website/internal/core/service/size_service"
	service "Eccomerce-website/internal/core/service/user_service"

	"Eccomerce-website/internal/infra/config"
	"Eccomerce-website/internal/infra/middleware"
	"Eccomerce-website/internal/infra/repository"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	// cwd, _ := os.Getwd()
	// fmt.Println("cwd :", cwd)
	errorMiddleware := middleware.ErrorMiddleware
	instance := gin.New()

	instance.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// instance.MaxMultipartMemory = 8 << 20 // 8MB maximum
	instance.Use(gin.Recovery())
	instance.Use(gin.Logger())
	instance.Use(errorMiddleware())
	v := validator.New()
	instance.Use(func(c *gin.Context) {
		c.Set("validator", v)
		c.Next()
	})

	engine := router.NewRouter(instance)

	dbConfig, httpServerConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDatabase(*dbConfig)
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
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(engine, userService)
	userController.InitRouter()

	// product category service
	categoryRepo := repository.NewProductCategoryRepository(db)
	categoryService := categoryservice.NewProductCategoryRepository(categoryRepo)
	categoryController := controller.NewCategoryController(engine, categoryService)
	categoryController.InitCategoryRouter()

	// color service
	colorRepo := repository.NewColorRepository(db)
	colorService := colorservice.NewColorService(colorRepo)
	colorController := controller.NewColorController(engine, colorService)
	colorController.InitColorRouter()

	// product service
	productRepo := repository.NewProductRepository(db)
	productService := productservice.NewProductService(productRepo)
	productController := controller.NewProductController(engine, productService)
	productController.InitProductRouter()

	// product item service
	productItemRepo := repository.NewProductItemRepository(db)
	productItemService := productitemservice.NewProductItemService(productItemRepo)
	productItemController := controller.NewProductItemController(engine, productItemService)
	productItemController.InitProductItemRouter()

	// product image service
	imageRepo := repository.NewProductImageRepository(db)
	productImageService := productimageservice.NewProductImageService(imageRepo)
	productImageController := controller.NewProductImageController(engine, productImageService)
	productImageController.InitProductImageRouter()

	// cart service
	cartRepo := repository.NewCartRepository(db)
	cartService := cartservice.NewCartService(cartRepo)
	cartController := controller.NewCartController(engine, cartService)
	cartController.InitCartRouter()

	// size service
	sizeRepo := repository.NewSizeRepository(db)
	sizeService := sizeservice.NewSizeService(sizeRepo)
	sizeController := controller.NewSizeController(engine, sizeService)
	sizeController.InitSizeRouter()

	// products service
	productsRepo := repository.NewProductsRepository(db)
	productsService := productsservice.NewProductsService(productsRepo)
	productsController := controller.NewProductsController(engine, productsService)
	productsController.InitProductsRouter()

	if err := server.Start(instance, *httpServerConfig); err != nil {
		log.Fatal(err)
	}
}
