package config

type Config struct {
	App *AppConfig
	DB  *DBConfig
	//Kafka *KafkaConfig
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type KafkaConfig struct {
	BrokerAddr string
	GroupID    string
	Topic      string
}
