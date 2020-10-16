package configuration

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Database DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Username      string
	Password      string
	Host          string
	Port          int
	Database      string
	ConnectionUri string
	Provider      string
}

func LoadConfiguration(configName string, configLocation string) (*Configuration, error) {
	configName = strings.TrimSuffix(configName, ".yml")
	configName = strings.TrimSuffix(configName, ".yaml")

	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%s", err)
		return nil, err
	}

	var configuration Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return nil, err
	}

	ValidateConfig(configuration)

	return &configuration, nil
}

func ValidateConfig(configuration Configuration) {
	dbConfig := configuration.Database
	if dbConfig.Provider == "" {
		panic("DB provider missing on configuration")
	}

	if dbConfig.ConnectionUri != "" {
		return
	}
	if dbConfig.Username != "" && dbConfig.Password != "" && dbConfig.Host != "" && dbConfig.Port != 0 && dbConfig.Database != "" {
		return
	}
	panic("Missing DB configuration")
}
