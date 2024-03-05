package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
	"github.com/mrvin/tasks-go/books/internal/storage"
	"google.golang.org/grpc"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv *grpc.Server
	conn net.Listener
	st   storage.Storage
	addr string
}

func New(conf *Conf, storage storage.Storage) (*Server, error) {
	var server Server

	server.st = storage

	var err error
	server.addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server.conn, err = net.Listen("tcp", server.addr)
	if err != nil {
		return nil, fmt.Errorf("establish tcp connection: %w", err)
	}

	var opts []grpc.ServerOption
	server.serv = grpc.NewServer(opts...)
	booksapi.RegisterBookServiceServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	slog.Info("Start gRPC server: " + s.addr)
	if err := s.serv.Serve(s.conn); err != nil {
		return fmt.Errorf("start grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	slog.Info("Stop gRPC server")
	s.serv.GracefulStop()
	s.conn.Close()
}
