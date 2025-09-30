package config

type Config struct {
	App     *AppConfig
	DB      *DBConfig
	Retries *DBRetriesConfig
	Kafka   *KafkaConfig
	Redis   *RedisConfig
	Cache   *Cache
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"open_conns"`
	MaxIdleConns int    `mapstructure:"idle_conns"`
	MaxLifetime  int    `mapstructure:"lifetime"`
}

type DBRetriesConfig struct {
	MaxRetries int `mapstructure:"db_retries"`
	RetryDelay int `mapstructure:"db_retry_delay"`
}

type KafkaConfig struct {
	BrokerAddr string `mapstructure:"broker_addr"`
	GroupID    string `mapstructure:"group_id"`
	Topic      string `mapstructure:"topic"`
	DLQTopic   string `mapstructure:"dlq_topic"`
	MaxRetries int    `mapstructure:"max_retries"`
}

type RedisConfig struct {
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
	User       string `mapstructure:"user"`
	DB         int    `mapstructure:"db"`
	MaxRetries int    `mapstructure:"max_retries"`
}
type Cache struct {
	CacheStartUpLimit int `mapstructure:"limit"`
}
