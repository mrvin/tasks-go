package config

import (
	"log/slog"
	"os"
)

type Config struct {
	ServerHTTPAddr string
	RedisCacheAddr string
	DnmonsterAddr  string
}

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if serverHTTPAddr := os.Getenv("SERVER_HTTP_ADDR"); serverHTTPAddr != "" {
		c.ServerHTTPAddr = serverHTTPAddr
	} else {
		slog.Warn("Empty server http addr")
	}
	if redisCacheAddr := os.Getenv("REDIS_CACHE_ADDR"); redisCacheAddr != "" {
		c.RedisCacheAddr = redisCacheAddr
	} else {
		slog.Warn("Empty redis cache addr")
	}
	if dnmonsterAddr := os.Getenv("DNMONSTER_ADDR"); dnmonsterAddr != "" {
		c.DnmonsterAddr = dnmonsterAddr
	} else {
		slog.Warn("Empty dnmonster addr")
	}
}
