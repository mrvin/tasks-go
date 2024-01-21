package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	createperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/create"
	deleteperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/delete"
	getperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/get"
	listpersons "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/list"
	updateperson "github.com/mrvin/tasks-go/persons/internal/httpserver/handlers/person/update"
	"github.com/mrvin/tasks-go/persons/internal/storage"
	"github.com/mrvin/tasks-go/persons/pkg/http/logger"
	"github.com/mrvin/tasks-go/persons/pkg/http/resolver"
	pathresolver "github.com/mrvin/tasks-go/persons/pkg/http/resolver/path"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.PersonStorage) *Server {
	res := pathresolver.New()

	res.Add(http.MethodPost+" /person", createperson.New(st))
	res.Add(http.MethodGet+" /person/", getperson.New(st))
	res.Add(http.MethodPut+" /person/", updateperson.New(st))
	res.Add(http.MethodDelete+" /person/", deleteperson.New(st))
	res.Add(http.MethodGet+" /list-persons", listpersons.New(st))

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
