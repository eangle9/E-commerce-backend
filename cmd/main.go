package main

import (
	// "Eccomerce-website/internal/controller"
	// "Eccomerce-website/internal/core/common/router"
	"Eccomerce-website/internal/controller"
	"Eccomerce-website/internal/core/common/router"
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

	userRepo := repository.NewUserRepository(db)
	// userService := service.NewUserService(userRepo)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(engine, userService)
	userController.InitRouter()

	if err := server.Start(instance, *httpServerConfig); err != nil {
		log.Fatal(err)
	}
}
