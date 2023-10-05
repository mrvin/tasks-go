package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/url-shortener/internal/app"
	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver/handler"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/logger"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/resolver"
	pathresolver "github.com/mrvin/tasks-go/url-shortener/pkg/http/resolver/path"
)

//nolint:tagliatelle
type ConfHTTPS struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

//nolint:tagliatelle
type Conf struct {
	Host    string    `yaml:"host"`
	Port    int       `yaml:"port"`
	IsHTTPS bool      `yaml:"is_https"`
	HTTPS   ConfHTTPS `yaml:"https"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, app *app.App) *Server {
	res := pathresolver.New()

	h := handler.New(app, res)

	res.Add("GET /health", h.HealthLivenessHandler)
	res.Add("POST /url", h.PutURL)
	res.Add("DELETE /url/", h.DeleteURL)

	loggerServer := logger.Logger{Inner: &Router{res}}

	return &Server{
		//nolint:exhaustivestruct,exhaustruct
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

func (s *Server) StartTLS(conf *ConfHTTPS) error {
	slog.Info("Start http server: https://" + s.Addr)
	if err := s.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil {
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
