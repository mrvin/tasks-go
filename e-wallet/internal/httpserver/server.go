package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	balancewallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/balance"
	createwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/create"
	historywallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/history"
	sendwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/send"

	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	"github.com/mrvin/tasks-go/e-wallet/pkg/http/logger"
	"github.com/mrvin/tasks-go/e-wallet/pkg/http/resolver"
	regexresolver "github.com/mrvin/tasks-go/e-wallet/pkg/http/resolver/regex"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.WalletStorage) *Server {
	res := regexresolver.New()

	regexpUUID := "([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})"
	res.Add(http.MethodPost+" /api/v1/wallet$", createwallet.New(st))
	res.Add(http.MethodPost+" /api/v1/wallet/"+regexpUUID+"/send", sendwallet.New(st))
	res.Add(http.MethodGet+" /api/v1/wallet/"+regexpUUID+"/history", historywallet.New(st))
	res.Add(http.MethodGet+" /api/v1/wallet/"+regexpUUID, balancewallet.New(st))

	loggerServer := logger.Logger{Inner: &Router{res}}

	return &Server{
		http.Server{
			Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Handler:      &loggerServer,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Minute,
		},
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

type Router struct {
	resolver.Resolver
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	if handlerFunc := r.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
	}

	http.NotFound(res, req)
}
