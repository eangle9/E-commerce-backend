package initiator

import (
	"Eccomerce-website/initiator/foundation"
	"fmt"
	"os"
)

//	@title			E-commerce API
//	@version		1.0
//	@description	This is a sample server for an Ecommerce platform.

//	@contact.name	Engdawork Yismaw
//	@contact.email	engdaworkyismaw9@gmail.com

//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization
//  @securityDefinitions.basic  BasicAuth
//	@host		localhost:9000
//	@schemes	http

func Initiate() {
	// Initiate Logger
	logger := foundation.InitLogger()
	log := logger.GetLogger()

	// Initiate Config
	log.Info("initializing config")
	configName := "config"
	if name := os.Getenv("CONFIG_NAME"); name != "" {
		configName = name
		log.Info(fmt.Sprintf("config name is set to %v", configName))
	} else {
		log.Info("using default config name 'config'")
	}
	foundation.InitConfig(configName, "config", log)
	log.Info("config initialized")

}
