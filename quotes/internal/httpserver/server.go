package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/quotes/internal/httpserver/handlers/health"
	createquote "github.com/mrvin/tasks-go/quotes/internal/httpserver/handlers/quote/create"
	deletequote "github.com/mrvin/tasks-go/quotes/internal/httpserver/handlers/quote/deleteq"
	getquote "github.com/mrvin/tasks-go/quotes/internal/httpserver/handlers/quote/get"
	listquotes "github.com/mrvin/tasks-go/quotes/internal/httpserver/handlers/quote/list"
	"github.com/mrvin/tasks-go/quotes/internal/storage"
	"github.com/mrvin/tasks-go/quotes/pkg/http/logger"
)

const (
	readTimeout  = 5  // in second
	writeTimeout = 10 // in second
	idleTimeout  = 1  // in minute
)

type Conf struct {
	Addr string
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.QuoteStorage) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", health.Health)

	mux.HandleFunc(http.MethodPost+" /quotes", createquote.New(st))
	mux.HandleFunc(http.MethodGet+" /quotes/random", getquote.New(st))
	mux.HandleFunc(http.MethodDelete+" /quotes/{id}", deletequote.New(st))

	mux.HandleFunc(http.MethodGet+" /quotes", listquotes.New(st))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustruct
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
			return
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
