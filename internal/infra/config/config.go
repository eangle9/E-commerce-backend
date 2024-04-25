package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Driver string
	Url    string
}

type HttpServerConfig struct {
	Port string
}

func LoadConfig() (*DatabaseConfig, *HttpServerConfig, error) {
	viper.SetConfigName("app")
	viper.AddConfigPath("../conf")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, err
	}

	dbConfig := DatabaseConfig{
		Driver: viper.GetString("database.driver"),
		Url:    viper.GetString("database.url"),
	}

	HttpServerConfig := HttpServerConfig{
		Port: fmt.Sprintf("%d", viper.GetInt("http_server.port")),
	}

	return &dbConfig, &HttpServerConfig, nil
}
