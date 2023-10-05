package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/app"
	"github.com/mrvin/tasks-go/url-shortener/internal/config"
	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver"
	"github.com/mrvin/tasks-go/url-shortener/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/url-shortener/internal/storage/sql"
)

type Config struct {
	HTTP   httpserver.Conf `yaml:"http"`
	DB     sqlstorage.Conf `yaml:"db"`
	Logger logger.Conf     `yaml:"logger"`
}

func main() {
	configFile := flag.String("config", "/etc/url-shortener/url-shortener.yml", "path to configuration file")
	flag.Parse()

	// init config
	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		stdlog.Printf("Parse config: %v", err)
		return
	}

	// init logger
	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v\n", err)
		return
	} else {
		slog.Info("Init logger", slog.String("level", conf.Logger.Level))
		defer func() {
			if err := logFile.Close(); err != nil {
				slog.Error("Close log file: " + err.Error())
			}
		}()
	}
	// init storage
	st, err := sqlstorage.New(context.Background(), &conf.DB)
	if err != nil {
		slog.Error("Failed to init storage: " + err.Error())
		return
	}
	slog.Info("Connected to database", slog.String("driver", conf.DB.Driver))

	// Start server
	server := httpserver.New(&conf.HTTP, app.New(st))

	if err := server.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
			return
		}
	}
}
