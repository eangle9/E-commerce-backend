package server

import (
	"Eccomerce-website/internal/infra/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Start(instance *gin.Engine, conf config.HttpServerConfig) error {
	// viper.SetConfigFile("../.env")
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.ReadInConfig()
	portKey := viper.Get("PORT")
	port := portKey.(string)
	// port := conf.Port
	// port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	err := instance.Run(":" + port)
	return err
}
