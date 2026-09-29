package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr      string
	Username  string
	Password  string
	DB        int
	UseTLS    bool
	tlsConfig *tls.Config
}

type Config struct {
	AppEnv       string
	HTTPAddr     string
	DatabaseURL  string
	Redis        RedisConfig
	KafkaBrokers string
}

func Load() (Config, error) {
	redisCfg, err := loadRedisConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppEnv:       getenv("APP_ENV", "local"),
		HTTPAddr:     listenAddr(),
		DatabaseURL:  getenv("DATABASE_URL", "postgres://ghaura:ghaura@localhost:5432/ghaura?sslmode=disable"),
		Redis:        redisCfg,
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:19092"),
	}, nil
}

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return normalizeListenAddr(port)
	}
	return getenv("HTTP_ADDR", ":8080")
}

func normalizeListenAddr(addr string) string {
	if strings.Contains(addr, ":") {
		return addr
	}
	return ":" + addr
}

func loadRedisConfig() (RedisConfig, error) {
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		return parseRedisURL(redisURL)
	}
	return RedisConfig{Addr: getenv("REDIS_ADDR", "localhost:6379")}, nil
}

func parseRedisURL(rawURL string) (RedisConfig, error) {
	opts, err := redis.ParseURL(rawURL)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("REDIS_URL: %w", err)
	}
	return RedisConfig{
		Addr:      opts.Addr,
		Username:  opts.Username,
		Password:  opts.Password,
		DB:        opts.DB,
		UseTLS:    opts.TLSConfig != nil,
		tlsConfig: opts.TLSConfig,
	}, nil
}

func (r RedisConfig) ClientOptions() *redis.Options {
	opts := &redis.Options{
		Addr:         r.Addr,
		Username:     r.Username,
		Password:     r.Password,
		DB:           r.DB,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		TLSConfig:    r.tlsConfig,
	}
	if opts.TLSConfig == nil && r.UseTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return opts
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
