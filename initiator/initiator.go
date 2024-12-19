package initiator

import (
	"Eccomerce-website/initiator/foundation"
	"Eccomerce-website/initiator/platform"
	"context"
	"fmt"
	"os"

	"github.com/eangle9/log"
	"github.com/spf13/viper"
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

func Initiate() {
	// Initiate Logger
	log := log.New(platform.InitLogger(), log.Options{})
	log.Info(context.Background(), "initialized logger")

	// Initiate Config
	log.Info(context.Background(), "initializing config")
	configName := "config"
	if name := os.Getenv("CONFIG_NAME"); name != "" {
		configName = name
		log.Info(context.Background(), fmt.Sprintf("config name is set to %v", configName))
	} else {
		log.Info(context.Background(), "using default config name 'config'")
	}
	foundation.InitConfig(configName, "config", log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing state")
	state := foundation.InitState(log)
	log.Info(context.Background(), "state initialized")

	// Initiate Database connection using pgx
	log.Info(context.Background(), "initializing database")
	conn := foundation.InitDB(viper.GetString("database.url"), log)
	log.Info(context.Background(), "database initialized")
	if viper.GetBool("migration.active") {
		log.Info(context.Background(), "initializing migration")
		m := foundation.InitiateMigration(viper.GetString("migration.path"), viper.GetString("database.url"), log)
		foundation.UpMigration(m, log)
		log.Info(context.Background(), "migration initialized")
	}
	fmt.Printf("state: %v\n", state)
	fmt.Printf("conn: %v\n", conn)
}
