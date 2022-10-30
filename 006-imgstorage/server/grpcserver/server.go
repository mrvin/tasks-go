package grpcserver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/mrvin/tasks-go/006-imgstorage/internal/imgstorageapi"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// numGoroutinUploadDownload - number of goroutines upload/download files.
const numGoroutinUploadDownload = 10

// numGoroutinList - number of goroutines get list files.
const numGoroutinList = 100

const chunkSize = 64 * 1024 // 64 KiB

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

func (s *Server) UploadImage(stream imgstorageapi.ImgStorage_UploadImageServer) error { //nolint: nosnakecase
	s.semaUploadDownload <- struct{}{}        // acquire token
	defer func() { <-s.semaUploadDownload }() // release token

	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("receive image name: %w", err)
	}

	name := req.GetName()

	fileImg, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	defer fileImg.Close()

	imageSize := 0
	for {
		if err := stream.Context().Err(); err != nil {
			return fmt.Errorf("termination due to context: %w", err)
		}
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("receive chunk data: %w", err)
		}

		chunk := req.GetChunkData()

		size, err := fileImg.Write(chunk)
		if err != nil {
			return fmt.Errorf("write chunk data: %w", err)
		}
		imageSize += size
	}

	res := &imgstorageapi.UploadImageResponse{
		Size: uint32(imageSize),
	}

	if err := stream.SendAndClose(res); err != nil {
		return fmt.Errorf("send response and close stream: %w", err)
	}

	log.Printf("Image \"%s\" upload, %d bytes", name, imageSize)

	return nil
}

func (s *Server) DownloadImage(nameImg *imgstorageapi.DownloadImageRequest, stream imgstorageapi.ImgStorage_DownloadImageServer) error { //nolint: nosnakecase
	s.semaUploadDownload <- struct{}{}        // acquire token
	defer func() { <-s.semaUploadDownload }() // release token

	name := nameImg.GetName()

	fileImage, err := os.Open(filepath.Join(s.dir, name))
	if err != nil {
		return fmt.Errorf("open image file: %w", err)
	}
	defer fileImage.Close()

	reader := bufio.NewReader(fileImage)
	buffer := make([]byte, chunkSize)

	imageSize := 0
	for {
		if err := stream.Context().Err(); err != nil {
			return fmt.Errorf("termination due to context: %w", err)
		}
		size, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read chunk to buffer: %w", err)
		}

		res := &imgstorageapi.DownloadImageResponse{
			ChunkData: buffer[:size],
		}

		if err := stream.Send(res); err != nil {
			return fmt.Errorf("send chunk: %w", err)
		}
		imageSize += size
	}

	log.Printf("Image \"%s\" download, %d bytes", name, imageSize)

	return nil
}

func (s *Server) GetListImage(ctx context.Context, _ *emptypb.Empty) (*imgstorageapi.GetListImageResponse, error) {
	s.semaList <- struct{}{}        // acquire token
	defer func() { <-s.semaList }() // release token

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("reading a image directory: %w", err)
	}

	listImg := make([]*imgstorageapi.ImageInfo, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return nil, fmt.Errorf("get information: %w", err)
			}
			listImg = append(listImg, &imgstorageapi.ImageInfo{Name: entry.Name(), ModifiedAt: timestamppb.New(info.ModTime())})
		}
	}

	log.Println("Get image list")

	return &imgstorageapi.GetListImageResponse{ImageInfo: listImg}, nil
}
