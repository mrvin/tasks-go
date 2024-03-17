package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/identidock/internal/config"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver/handlers/getidenticon"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver/handlers/mainform"
	"github.com/mrvin/tasks-go/identidock/pkg/http/logger"
)

type Server struct {
	http.Server
}

func New(conf *config.Config, cache *redis.Client) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /main", mainform.NewGet())
	mux.HandleFunc(http.MethodPost+" /main", mainform.NewPost())
	mux.HandleFunc(http.MethodGet+" /monster/", getidenticon.New(conf.DnmonsterAddr, cache))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustruct
		http.Server{
			Addr:         conf.ServerHTTPAddr,
			Handler:      &loggerServer,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Minute,
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
