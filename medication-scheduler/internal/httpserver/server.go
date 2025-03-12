package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/httpserver/handlers"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
	"github.com/mrvin/tasks-go/medication-scheduler/pkg/http/logger"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

type Conf struct {
	PeriodNextTakings time.Duration
	Addr              string
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.SchedulerStorage) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", handlers.Health)

	mux.HandleFunc(http.MethodPost+" /schedule", handlers.NewCreateSchedule(st))
	mux.HandleFunc(http.MethodGet+" /schedule", handlers.NewGetSchedule(st))

	mux.HandleFunc(http.MethodGet+" /schedules", handlers.NewListSchedulesIDs(st))
	mux.HandleFunc(http.MethodGet+" /next_takings", handlers.NewGetNextTakings(st, conf.PeriodNextTakings))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustivestruct,exhaustruct
		http.Server{
			Addr:         conf.Addr,
			Handler:      &loggerServer,
			ReadTimeout:  readTimeout * time.Second,
			WriteTimeout: writeTimeout * time.Second,
			IdleTimeout:  idleTimeout * time.Minute,
		},
	}
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
		}
	}()
	slog.Info("Start http server: http://" + s.Addr)

	<-ctx.Done()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Failed to stop http server: " + err.Error())
		return
	}
	slog.Info("Stop http server")
}
