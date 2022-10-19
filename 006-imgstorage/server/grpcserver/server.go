package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/mrvin/tasks-go/006-imgstorage/internal/imgstorageapi"
	"google.golang.org/grpc"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Dir  string `yaml:"dir"`
}

type Server struct {
	serv *grpc.Server
	ln   net.Listener
	addr string
	dir  string
}

func New(conf *Config) (*Server, error) {
	var server Server

	server.dir = conf.Dir

	var err error
	server.addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server.ln, err = net.Listen("tcp", server.addr)
	if err != nil {
		return nil, fmt.Errorf("establish tcp connection: %w", err)
	}
	server.serv = grpc.NewServer()
	imgstorageapi.RegisterImgStorageServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	log.Printf("Start gRPC server: %s", s.addr)
	if err := s.serv.Serve(s.ln); err != nil {
		return fmt.Errorf("start gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	log.Print("Stop gRPC server")
	s.serv.GracefulStop()
	s.ln.Close()
}

func (s *Server) UploadImg(ctx context.Context, img *imgstorageapi.Img) (*imgstorageapi.Null, error) {
	name := img.Name
	image := img.Img

	fileImg, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		log.Printf("Create image: %v", err)
		return nil, err
	}

	size, err := fileImg.Write(image)

	if closeErr := fileImg.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		log.Printf("Image \"%s\" upload, %d bytes", name, size)
	}

	return &imgstorageapi.Null{}, err
}

func (s *Server) DownloadImg(ctx context.Context, nameImg *imgstorageapi.NameImg) (*imgstorageapi.Img, error) {
	name := nameImg.Name

	image, err := os.ReadFile(filepath.Join(s.dir, name))
	if err != nil {
		log.Printf("Read image: %v", err)
		return nil, err
	}

	log.Printf("Image \"%s\" download, %d bytes", name, len(image))

	return &imgstorageapi.Img{Name: name, Img: image}, nil
}
