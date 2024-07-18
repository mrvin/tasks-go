package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	clienthttp "github.com/mrvin/tasks-go/hh-client-go/internal/client/http"
	"github.com/mrvin/tasks-go/hh-client-go/internal/config"
	"github.com/mrvin/tasks-go/hh-client-go/internal/logger"
)

type Config struct {
	AuthHH clienthttp.ConfAPIhh `yaml:"auth_hh"`
	HTTP   clienthttp.Conf      `yaml:"http"`
	Logger logger.Conf          `yaml:"logger"`
}

func main() {
	configFile := flag.String("config", "/etc/hh-client-go/hh-client-go.yml", "path to configuration file")
	flag.Parse()

	appInfo := clienthttp.AppInfo{
		Name:    "hh-client-go",
		Version: "1.0",
		Email:   "v.v.vinogradovv@gmail.com",
	}

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
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()
	slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx,
		syscall.SIGINT, /*(Control-C)*/
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	client, err := clienthttp.New(ctx, &conf.HTTP, &conf.AuthHH, &appInfo)
	if err != nil {
		slog.Error("Create http client: " + err.Error())
		return
	}

	slResumeID, err := client.ListResumeID(ctx)
	if err != nil {
		slog.Error("Get list resume ID: " + err.Error())
		return
	}

	chDoneAutoUpdateResumes := make(chan struct{})
	funcAutoUpdateResumes := func() {
		client.AutoUpdateResumes(ctx, slResumeID, chDoneAutoUpdateResumes)
	}

	durationBeforeStart := durationBeforeStartAutoUpdateResumes()
	slog.Info("Auto update will start", slog.Duration("duration", durationBeforeStart))
	time.AfterFunc(durationBeforeStart, funcAutoUpdateResumes)

	<-chDoneAutoUpdateResumes
	close(chDoneAutoUpdateResumes)
}

func durationBeforeStartAutoUpdateResumes() time.Duration {
	var timeStartAutoUpdateResumes time.Time

	timeNow := time.Now()
	year, month, day := timeNow.Date()

	hour := timeNow.Hour()
	switch {
	case hour < 2:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 2, 0, 0, 0, timeNow.Location())
	case hour < 6:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 6, 0, 0, 0, timeNow.Location())
	case hour < 10:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 10, 0, 0, 0, timeNow.Location())
	case hour < 14:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 14, 0, 0, 0, timeNow.Location())
	case hour < 18:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 18, 0, 0, 0, timeNow.Location())
	case hour < 22:
		timeStartAutoUpdateResumes = time.Date(year, month, day, 22, 0, 0, 0, timeNow.Location())
	case hour < 24:
		timeStartAutoUpdateResumes = time.Date(year, month, day+1, 2, 0, 0, 0, timeNow.Location())
	}

	return timeStartAutoUpdateResumes.Sub(timeNow)
}
