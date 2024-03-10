package httpserver

import (
	"context"

	"fmt"

	"log/slog"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver/handlers/getidenticon"
	"github.com/mrvin/tasks-go/identidock/internal/httpserver/handlers/mainform"
	"github.com/mrvin/tasks-go/identidock/pkg/http/logger"
)

type Config struct {
	Addr          string
	DnmonsterAddr string
}

type Server struct {
	http.Server
	Addr string
}

func New(serverHTTPAddr, dnmonsterAddr string, cache *redis.Client) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /main", mainform.NewGet())
	mux.HandleFunc(http.MethodPost+" /main", mainform.NewPost())
	mux.HandleFunc(http.MethodGet+" /monster/", getidenticon.New(dnmonsterAddr, cache))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		http.Server{
			Addr:         serverHTTPAddr,
			Handler:      &loggerServer,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Minute,
		},
		serverHTTPAddr,
	}
}

func (s *Server) Start() error {
	slog.Info("Start http server: http://" + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stop http server")
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
