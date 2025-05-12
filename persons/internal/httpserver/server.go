package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	// import docs generated with swag init.
	_ "github.com/mrvin/tasks-go/persons/docs"
	"github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/health"
	createperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/create"
	deleteperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/deletep"
	getperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/get"
	listpersons "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/list"
	updateperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/update"
	updatefullperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/updatefull"
	"github.com/mrvin/tasks-go/persons/internal/storage"
	"github.com/mrvin/tasks-go/persons/pkg/http/logger"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

// @host		localhost:8080
func New(conf *Conf, st storage.PersonStorage) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodGet+" /health", health.Health)

	mux.HandleFunc(http.MethodPost+" /persons", createperson.New(st))
	mux.HandleFunc(http.MethodGet+" /persons/{id}", getperson.New(st))
	mux.HandleFunc(http.MethodPut+" /persons/{id}", updatefullperson.New(st))
	mux.HandleFunc(http.MethodPatch+" /persons/{id}", updateperson.New(st))
	mux.HandleFunc(http.MethodDelete+" /persons/{id}", deleteperson.New(st))

	mux.HandleFunc(http.MethodGet+" /persons", listpersons.New(st))

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

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
