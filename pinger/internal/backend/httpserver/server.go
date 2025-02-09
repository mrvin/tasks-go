package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/pinger/internal/backend/httpserver/handlers"
	"github.com/mrvin/tasks-go/pinger/internal/storage"
	"github.com/mrvin/tasks-go/pinger/pkg/http/logger"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

type Conf struct {
	Addr string
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.PingerStorage) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", handlers.Health)

	mux.HandleFunc(http.MethodPost+" /hosts", handlers.NewCreateHost(st))
	mux.HandleFunc(http.MethodGet+" /hosts", handlers.NewListHost(st))

	mux.HandleFunc(http.MethodPost+" /pings", handlers.NewCreatePing(st))
	mux.HandleFunc(http.MethodGet+" /pings", handlers.NewListLatestPing(st))

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
