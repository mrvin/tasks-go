package main

import (
	"context"
	stdlog "log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/pinger/internal/logger"
	"github.com/mrvin/tasks-go/pinger/internal/pinger/app"
	"github.com/mrvin/tasks-go/pinger/internal/pinger/config"
)

func main() {
	// init config
	var conf config.Config
	conf.LoadFromEnv()

	// init logger
	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v\n", err)
		return
	}
	slog.Info("Init logger", slog.String("level", conf.Logger.Level))
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	app.Run(ctx, &conf.App, &conf.HTTP)
}
