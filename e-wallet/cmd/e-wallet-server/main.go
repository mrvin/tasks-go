package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/e-wallet/internal/app"
	"github.com/mrvin/tasks-go/e-wallet/internal/config"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver"
	"github.com/mrvin/tasks-go/e-wallet/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
)

type Config struct {
	App    app.Conf        `yaml:"app"`
	DB     sqlstorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf `yaml:"http"`
	Logger logger.Conf     `yaml:"logger"`
}

func main() {
	configFile := flag.String("config", "/etc/e-wallet/e-wallet.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Parse config: %v", err)
		return
	}

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

	serverHTTP := httpserver.New(&conf.App, &conf.HTTP, storage)

	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	serverHTTP.Run(ctx)
}
