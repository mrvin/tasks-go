package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/mrvin/tasks-go/photo-gallery/internal/config"
	"github.com/mrvin/tasks-go/photo-gallery/internal/httpserver"
	"github.com/mrvin/tasks-go/photo-gallery/internal/logger"
	sqlitestorage "github.com/mrvin/tasks-go/photo-gallery/internal/storage/sqlite"
)

type Config struct {
	DB     sqlitestorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf    `yaml:"http"`
	Logger logger.Conf        `yaml:"logger"`
}

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/photo-gallery/photo-gallery.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Parse config: %v", err)
		return
	}

	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		log.Printf("Init logger: %v", err)
		return
	}
	slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()

	slog.Info("Storage in sql database")
	storage, err := sqlitestorage.New(ctx, &conf.DB)
	if err != nil {
		slog.Error("Failed to init storage: " + err.Error())
		return
	}
	defer func() {
		if err := storage.Close(); err != nil {
			slog.Error("Failed to close storage: " + err.Error())
		} else {
			slog.Info("Closing the database connection")
		}
	}()

	slog.Info("Connected to database")

	if err := os.MkdirAll(conf.HTTP.DirPhotos, 0750); err != nil {
		slog.Error("Failed to create dir: " + err.Error())
		return
	}

	serverHTTP := httpserver.New(&conf.HTTP, storage)

	err = serverHTTP.Start()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start http server: " + err.Error())
		return
	}
}
