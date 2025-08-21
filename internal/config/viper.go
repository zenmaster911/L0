package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func MustLoad() (*Config, *KafkaConfig) {
	configPath := "local.yaml"
	viper.SetConfigFile(configPath)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env file: %s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error in reading config: %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error in unmarshalign config: %s", err)
	}
	var kafka KafkaConfig
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	kafka.BrokerAddr = os.Getenv("KAFKA_BROKER_ADDR")
	kafka.GroupID = os.Getenv("GROUP_ID")
	kafka.Topic = os.Getenv("TOPIC_NAME")

	return &cfg, &kafka
}
