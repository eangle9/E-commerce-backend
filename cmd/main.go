package main

import (
	// "Eccomerce-website/internal/controller"
	// "Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/controller"
	"Eccomerce-website/internal/core/common/router"
	categoryservice "Eccomerce-website/internal/core/service/category_service"
	colorservice "Eccomerce-website/internal/core/service/color_service"
	productitemservice "Eccomerce-website/internal/core/service/product_item_service"
	productservice "Eccomerce-website/internal/core/service/product_service"
	service "Eccomerce-website/internal/core/service/user_service"

	// "Eccomerce-website/internal/core/service"

	// "Eccomerce-website/internal/core/port/service"
	"Eccomerce-website/internal/core/server"

	// "Eccomerce-website/internal/core/service"
	"Eccomerce-website/internal/infra/config"
	"Eccomerce-website/internal/infra/middleware"
	"Eccomerce-website/internal/infra/repository"

	// "Eccomerce-website/schema"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/go-playground/validator/v10"
)

func main() {
	// cwd, _ := os.Getwd()
	// fmt.Println("cwd :", cwd)
	errorMiddleware := middleware.ErrorMiddleware
	instance := gin.New()
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

	if err := server.Start(instance, *httpServerConfig); err != nil {
		log.Fatal(err)
	}
}
