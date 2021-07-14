package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server       ServerConfigurations
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
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}
