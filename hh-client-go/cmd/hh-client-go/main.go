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

//nolint:tagliatelle
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
	slog.Info("Auto update will start", slog.String("duration", durationBeforeStart.String()))
	time.AfterFunc(durationBeforeStart, funcAutoUpdateResumes)

	<-chDoneAutoUpdateResumes
	close(chDoneAutoUpdateResumes)
}

func durationBeforeStartAutoUpdateResumes() time.Duration {
	var timeStartAutoUpdateResumes time.Time

	timeNow := time.Now()
	year, month, day := timeNow.Date()

	hour := timeNow.Hour()
	hours := [...]int{6, 10, 14, 18, 22}

	for _, h := range hours {
		if hour < h {
			timeStartAutoUpdateResumes = time.Date(year, month, day, h, 0, 0, 0, timeNow.Location())
			break
		}
	}
	if hour >= hours[len(hours)-1] {
		timeStartAutoUpdateResumes = time.Date(year, month, day+1, hours[0], 0, 0, 0, timeNow.Location())
	}

	return timeStartAutoUpdateResumes.Sub(timeNow)
}
