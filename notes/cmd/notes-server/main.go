package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/notes/internal/config"
	"github.com/mrvin/tasks-go/notes/internal/httpserver"
	"github.com/mrvin/tasks-go/notes/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/notes/internal/storage/sql"
)

func main() {
	// init config
	var conf config.Config
	conf.LoadFromEnv()

	// init logger
	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		log.Printf("Init logger: %v\n", err)
		return
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()
	slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))

	// init storage
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

	// Start server
	server := httpserver.New(&conf.HTTP, storage)

	if err := server.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
			return
		}
	}
}
