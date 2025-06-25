package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	creategood "github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/good/create"
	deletegood "github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/good/delete"
	listgoods "github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/good/list"
	reprioritizegood "github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/good/reprioritize"
	updategood "github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/good/update"
	"github.com/mrvin/tasks-go/goods/internal/httpserver/handlers/health"
	natsmq "github.com/mrvin/tasks-go/goods/internal/queue/nats"
	"github.com/mrvin/tasks-go/goods/internal/storage"
	"github.com/mrvin/tasks-go/goods/pkg/http/logger"
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

func New(conf *Conf, st storage.GoodsStorage, mq *natsmq.Queue) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", health.Health)

	mux.HandleFunc(http.MethodGet+" /goods/list", listgoods.New(st))
	mux.HandleFunc(http.MethodPost+" /good/create", creategood.New(st, mq))
	mux.HandleFunc(http.MethodPatch+" /good/update", updategood.New(st, mq))
	mux.HandleFunc(http.MethodPatch+" /good/reprioritize", reprioritizegood.New(st, mq))
	mux.HandleFunc(http.MethodDelete+" /good/remove", deletegood.New(st, mq))

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
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
			defer cancel()
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
