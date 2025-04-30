package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/persons/internal/config"
	"github.com/mrvin/tasks-go/persons/internal/httpserver"
	"github.com/mrvin/tasks-go/persons/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/persons/internal/storage/sql"
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

	ctx := context.Background()
	slog.Info("Storage in sql database")
	storage, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		slog.Error("Failed to init storage: " + err.Error())
		return
	}
	slog.Info("Connected to database")
	defer func() {
		if err := storage.Close(); err != nil {
			slog.Error("Failed to close storage: " + err.Error())
		} else {
			slog.Info("Closing the database connection")
		}
	}()

	server := httpserver.New(&conf.HTTP, storage)

	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	server.Run(ctx)
}
