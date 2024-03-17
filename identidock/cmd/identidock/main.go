package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/identidock/internal/config"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver"
)

func main() {
	var conf config.Config

	conf.LoadFromEnv()

	ctx := context.Background()
	cache := redis.NewClient(
		//nolint:exhaustruct
		&redis.Options{
			Addr:     conf.RedisCacheAddr,
			Password: "",
			DB:       0,
		})

	if err := cache.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	serverHTTP := httpserver.New(&conf, cache)

	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	serverHTTP.Run(ctx)
}
