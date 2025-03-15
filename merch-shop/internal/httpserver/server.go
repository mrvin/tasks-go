package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/merch-shop/internal/app"
	"github.com/mrvin/tasks-go/merch-shop/internal/auth"
	"github.com/mrvin/tasks-go/merch-shop/internal/httpserver/handlers"
	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
	"github.com/mrvin/tasks-go/merch-shop/pkg/http/logger"
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

func New(conf *Conf, appConf *app.Conf, st storage.ShopStorage, a *auth.AuthService) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", handlers.Health)
	mux.HandleFunc(http.MethodPost+" /api/auth", handlers.NewAuth(appConf, st, a))
	mux.HandleFunc(http.MethodPost+" /api/sendCoin", a.Auth(handlers.NewSendCoin(st)))

	mux.HandleFunc(http.MethodGet+" /api/buy/{productName}", a.Auth(handlers.NewBuyProduct(st)))
	mux.HandleFunc(http.MethodGet+" /api/info", a.Auth(handlers.NewInfo(st)))

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
