package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/mrvin/tasks-go/e-wallet/internal/config"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver"
	"github.com/mrvin/tasks-go/e-wallet/internal/logger"

	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
)

type Config struct {
	DB     sqlstorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf `yaml:"http"`
	Logger logger.Conf     `yaml:"logger"`
}

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/e-wallet/e-wallet.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		stdlog.Printf("Parse config: %v", err)
		return
	}

	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v\n", err)
		return
	} else {
		slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))
		defer func() {
			if err := logFile.Close(); err != nil {
				slog.Error("Close log file: " + err.Error())
			}
		}()
	}

	slog.Info("Storage in sql database")
	storage, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		slog.Error("Failed to init storage: " + err.Error())
		return
	}
	slog.Info("Connected to database")

	serverHTTP := httpserver.New(&conf.HTTP, storage)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		err = serverHTTP.Start()
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
			return
		}
	}()

	wg.Wait()
}
