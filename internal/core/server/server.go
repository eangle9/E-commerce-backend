package server

import (
	"Eccomerce-website/internal/infra/config"

	"github.com/gin-gonic/gin"
)

func Start(instance *gin.Engine, conf config.HttpServerConfig) error {
	port := conf.Port

	if port == "" {
		port = "8000"
	}

	err := instance.Run(":" + port)
	return err
}
