package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server ServerConfigurations
}

type ServerConfigurations struct {
	port                   int
	idleTimeout            int
	readTimeout            int
	timeoutContextDuration int
}

func SetConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("viper configuration reading error")
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("viper configuration unmarshalling error")
	}
}
