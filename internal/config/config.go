package config

type Config struct {
	App   *AppConfig
	DB    *DBConfig
	Kafka *KafkaConfig
	Redis *RedisConfig
	Cache *Cache
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
	BrokerAddr string `mapstructure:"broker_addr"`
	GroupID    string `mapstructure:"group_id"`
	Topic      string `mapstructure:"topic"`
}

type RedisConfig struct {
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
	User       string `mapstructure:"user"`
	DB         int    `mapstructure:"db"`
	MaxRetries int    `mapstructure:"max_retries"`
	// DialTimeout time.Duration `mapstructure:"dial_timeout"`
	// Timeout     time.Duration `mapstructure:"timeout"`
}
type Cache struct {
	CacheStartUpLimit int `mapstructure:"limit"`
}
