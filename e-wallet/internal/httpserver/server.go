package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	balancewallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/balance"
	createwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/create"
	depositwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/deposit"
	historywallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/history"
	sendwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/send"
	withdrawwallet "github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/withdraw"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	"github.com/mrvin/tasks-go/e-wallet/pkg/http/logger"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.WalletStorage) *Server {
	mux := http.NewServeMux()
	path := " /api/v1/wallet/"

	mux.HandleFunc(http.MethodPost+path[:len(path)-1], createwallet.New(st))
	mux.HandleFunc(http.MethodPost+path+"{walletID}/send", sendwallet.New(st))
	mux.HandleFunc(http.MethodPost+path+"{walletID}/deposit", depositwallet.New(st))
	mux.HandleFunc(http.MethodPost+path+"{walletID}/withdraw", withdrawwallet.New(st))

	mux.HandleFunc(http.MethodGet+path+"{walletID}/history", historywallet.New(st))
	mux.HandleFunc(http.MethodGet+path+"{walletID}", balancewallet.New(st))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustruct
		http.Server{
			Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
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
