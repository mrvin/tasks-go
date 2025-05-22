package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/quotes/internal/config"
	"github.com/mrvin/tasks-go/quotes/internal/httpserver"
	"github.com/mrvin/tasks-go/quotes/internal/logger"
	memorystorage "github.com/mrvin/tasks-go/quotes/internal/storage/memory"
)

func main() {
	// init config
	var conf config.Config
	conf.LoadFromEnv()

	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		log.Printf("Init logger: %v\n", err)
		return
	}
	slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()

	slog.Info("Storage in memory")
	storage := memorystorage.New()

	server := httpserver.New(&conf.HTTP, storage)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	server.Run(ctx)
}
