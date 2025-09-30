package config

import (
	"log"

	"github.com/spf13/viper"
)

func MustLoad() *Config {
	configPath := "local.yaml"
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error in reading config: %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error in unmarshalign config: %s", err)
	}

	return &cfg
}
