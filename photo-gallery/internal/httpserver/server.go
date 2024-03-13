package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/photo-gallery/internal/httpserver/handlers/photo/delete"
	"github.com/mrvin/tasks-go/photo-gallery/internal/httpserver/handlers/photo/list"
	"github.com/mrvin/tasks-go/photo-gallery/internal/httpserver/handlers/photo/save"
	"github.com/mrvin/tasks-go/photo-gallery/internal/storage"
	"github.com/mrvin/tasks-go/photo-gallery/pkg/http/logger"
)

type Conf struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	DirPhotos string `yaml:"dirPhotos"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, st storage.PhotoStorage) *Server {
	mux := http.NewServeMux()

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	path := "/api/v1/photo"

	mux.HandleFunc(http.MethodPost+" "+path, save.New(st, conf.DirPhotos, addr, path))
	mux.HandleFunc(http.MethodDelete+" "+path, delete.New(st, conf.DirPhotos))
	mux.HandleFunc(http.MethodGet+" /api/v1/listphotos", list.New(st))

	mux.Handle(path+"/", http.StripPrefix(path+"/", http.FileServer(http.Dir(conf.DirPhotos))))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		http.Server{
			Addr:         addr,
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
