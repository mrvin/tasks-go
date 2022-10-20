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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// numGoroutinUploadDownload - number of goroutines upload/download files.
const numGoroutinUploadDownload = 10

// numGoroutinList - number of goroutines get list files.
const numGoroutinList = 100

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Dir  string `yaml:"dir"`
}

type Server struct {
	serv               *grpc.Server
	ln                 net.Listener
	addr               string
	dir                string
	semaUploadDownload chan struct{}
	semaList           chan struct{}
}

func New(conf *Config) (*Server, error) {
	var server Server

	server.dir = conf.Dir
	server.semaUploadDownload = make(chan struct{}, numGoroutinUploadDownload)
	server.semaList = make(chan struct{}, numGoroutinList)

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
	s.semaUploadDownload <- struct{}{}        // acquire token
	defer func() { <-s.semaUploadDownload }() // release token

	name := img.Name
	image := img.Img

	fileImg, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		return nil, fmt.Errorf("create image: %w", err)
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
	s.semaUploadDownload <- struct{}{}        // acquire token
	defer func() { <-s.semaUploadDownload }() // release token

	name := nameImg.Name

	image, err := os.ReadFile(filepath.Join(s.dir, name))
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}

	log.Printf("Image \"%s\" download, %d bytes", name, len(image))

	return &imgstorageapi.Img{Name: name, Img: image}, nil
}

func (s *Server) GetListImg(ctx context.Context, _ *imgstorageapi.Null) (*imgstorageapi.ListImg, error) {
	s.semaList <- struct{}{}        // acquire token
	defer func() { <-s.semaList }() // release token

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("reading a image directory: %w", err)
	}

	listImg := make([]*imgstorageapi.InfImg, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return nil, fmt.Errorf("get information: %w", err)
			}
			listImg = append(listImg, &imgstorageapi.InfImg{Name: entry.Name(), ModifiedAt: timestamppb.New(info.ModTime())})
		}
	}

	log.Println("Get image list")

	return &imgstorageapi.ListImg{InfImg: listImg}, nil
}
