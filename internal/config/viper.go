package config

import (
	"fmt"
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
	fmt.Println(cfg.Cache.CacheStartUpLimit)
	//var kafka KafkaConfig
	// cfg.DB.Password = os.Getenv("DB_PASSWORD")
	// kafka.BrokerAddr = os.Getenv("KAFKA_BROKER_ADDR")
	// kafka.GroupID = os.Getenv("GROUP_ID")
	// kafka.Topic = os.Getenv("TOPIC_NAME")

	return &cfg
}
