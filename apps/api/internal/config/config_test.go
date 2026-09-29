package config

import (
	"strings"
	"testing"
)

func TestListenAddrPrefersPORTOverHTTPAddr(t *testing.T) {
	t.Setenv("PORT", "4096")
	t.Setenv("HTTP_ADDR", ":8080")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.HTTPAddr != ":4096" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":4096")
	}
}

func TestListenAddrKeepsHostPortFromPORT(t *testing.T) {
	t.Setenv("PORT", "0.0.0.0:9000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.HTTPAddr != "0.0.0.0:9000" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "0.0.0.0:9000")
	}
}

func TestListenAddrFallsBackToHTTPAddr(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("HTTP_ADDR", ":9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":9090")
	}
}

func TestListenAddrDefaultsTo8080(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("HTTP_ADDR", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
}

func TestRedisPrefersURLOverAddr(t *testing.T) {
	t.Setenv("REDIS_URL", "redis://default:s3cret@redis.railway.internal:6379/2")
	t.Setenv("REDIS_ADDR", "localhost:6379")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Redis.Addr != "redis.railway.internal:6379" {
		t.Fatalf("Redis.Addr = %q, want host:port from REDIS_URL", cfg.Redis.Addr)
	}
	if cfg.Redis.Username != "default" {
		t.Fatalf("Redis.Username = %q, want default", cfg.Redis.Username)
	}
	if cfg.Redis.Password != "s3cret" {
		t.Fatalf("Redis.Password = %q, want s3cret", cfg.Redis.Password)
	}
	if cfg.Redis.DB != 2 {
		t.Fatalf("Redis.DB = %d, want 2", cfg.Redis.DB)
	}
	if cfg.Redis.UseTLS {
		t.Fatal("expected redis:// to disable TLS")
	}
}

func TestRedisParsesRedissURL(t *testing.T) {
	t.Setenv("REDIS_URL", "rediss://:hunter2@secure.example:6380")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Redis.Addr != "secure.example:6380" {
		t.Fatalf("Redis.Addr = %q, want secure.example:6380", cfg.Redis.Addr)
	}
	if cfg.Redis.Password != "hunter2" {
		t.Fatalf("Redis.Password = %q, want hunter2", cfg.Redis.Password)
	}
	if !cfg.Redis.UseTLS {
		t.Fatal("expected rediss:// to enable TLS")
	}

	opts := cfg.Redis.ClientOptions()
	if opts.TLSConfig == nil {
		t.Fatal("ClientOptions TLSConfig is nil for rediss://")
	}
	if opts.Addr != cfg.Redis.Addr || opts.Password != cfg.Redis.Password {
		t.Fatalf("ClientOptions did not copy addr/password: %+v", opts)
	}
}

func TestRedisFallsBackToAddr(t *testing.T) {
	t.Setenv("REDIS_URL", "")
	t.Setenv("REDIS_ADDR", "redis:6379")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Redis.Addr != "redis:6379" {
		t.Fatalf("Redis.Addr = %q, want redis:6379", cfg.Redis.Addr)
	}
	if cfg.Redis.Password != "" || cfg.Redis.UseTLS {
		t.Fatalf("expected empty password and no TLS for REDIS_ADDR, got %+v", cfg.Redis)
	}
}

func TestRedisDefaultsToLocalhost(t *testing.T) {
	t.Setenv("REDIS_URL", "")
	t.Setenv("REDIS_ADDR", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Fatalf("Redis.Addr = %q, want localhost:6379", cfg.Redis.Addr)
	}
}

func TestRedisRejectsInvalidURL(t *testing.T) {
	t.Setenv("REDIS_URL", "not-a-redis-url")

	_, err := Load()
	if err == nil {
		t.Fatal("expected invalid REDIS_URL to fail")
	}
	if !strings.Contains(err.Error(), "REDIS_URL") {
		t.Fatalf("error %q should mention REDIS_URL", err)
	}
}
