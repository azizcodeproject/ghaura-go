package config

import "os"

type Config struct {
	AppEnv       string
	HTTPAddr     string
	DatabaseURL  string
	RedisAddr    string
	KafkaBrokers string
}

func Load() Config {
	return Config{
		AppEnv:       getenv("APP_ENV", "local"),
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:  getenv("DATABASE_URL", "postgres://ghaura:ghaura@localhost:5432/ghaura?sslmode=disable"),
		RedisAddr:    getenv("REDIS_ADDR", "localhost:6379"),
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:19092"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
