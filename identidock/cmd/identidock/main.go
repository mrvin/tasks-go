package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver"
)

type Config struct {
	ServerHTTPAddr string
	RedisCacheAddr string
	DnmonsterAddr  string
}

var ctx = context.Background()
var confDnmonsterStr string

func main() {
	conf := Config{
		ServerHTTPAddr: "identidock:8888",
		RedisCacheAddr: "redis:6379",
		DnmonsterAddr:  "dnmonster:8080",
	}

	conf.LoadFromEnv()

	cache := redis.NewClient(&redis.Options{
		Addr:     conf.RedisCacheAddr,
		Password: "",
		DB:       0,
	})

	if err := cache.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	serverHTTP := httpserver.New(conf.ServerHTTPAddr, conf.DnmonsterAddr, cache)

	if err := serverHTTP.Start(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start http server: " + err.Error())
		return
	}
}

// LoadFromEnv will load configuration solely from the environment
func (c *Config) LoadFromEnv() {
	if serverHTTPAddr := os.Getenv("SERVER_HTTP_ADDR"); serverHTTPAddr != "" {
		c.ServerHTTPAddr = serverHTTPAddr
	}
	if redisCacheAddr := os.Getenv("REDIS_CACHE_ADDR"); redisCacheAddr != "" {
		c.RedisCacheAddr = redisCacheAddr
	}
	if dnmonsterAddr := os.Getenv("DNMONSTER_ADDR"); dnmonsterAddr != "" {
		c.DnmonsterAddr = dnmonsterAddr
	}
}
