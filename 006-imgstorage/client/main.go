package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mrvin/tasks-go/006-imgstorage/internal/imgstorageapi"
	"google.golang.org/grpc"
)

type Client struct {
	client imgstorageapi.ImgStorageClient
}

var ctx = context.Background()

func main() {
	var port int
	var host string
	var fileName string
	var upload, download bool

	flag.IntVar(&port, "port", 55555, "port")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.StringVar(&fileName, "name", "", "file name")
	flag.BoolVar(&upload, "upload", false, "upload image")
	flag.BoolVar(&download, "download", false, "download image")
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Printf("gRPC connect: %v", err)
		return
	}
	defer conn.Close()

	var c Client
	c.client = imgstorageapi.NewImgStorageClient(conn)

	if upload && download {
		log.Printf("Can't upload and download")
		return
	}

	if upload {
		if err := c.uploadImg(fileName); err != nil {
			log.Printf("Upload image: %v", err)
		}
	}

	if download {
		if err := c.downloadImg(fileName); err != nil {
			log.Printf("Download image: %v", err)
		}
	}
}

func (c *Client) uploadImg(pathToFile string) error {
	image, err := os.ReadFile(pathToFile)
	if err != nil {
		return fmt.Errorf("read image: %w", err)
	}

	fileName := filepath.Base(pathToFile)

	req := &imgstorageapi.Img{Name: fileName, Img: image}

	_, err = c.client.UploadImg(ctx, req)
	if err != nil {
		return fmt.Errorf("upload image: %w", err)
	}

	log.Printf("Upload image \"%s\"", fileName)

	return nil
}

func (c *Client) downloadImg(fileName string) error {
	req := &imgstorageapi.NameImg{Name: fileName}

	pbImg, err := c.client.DownloadImg(ctx, req)
	if err != nil {
		return fmt.Errorf("download image: %w", err)
	}

	fileImg, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}

	size, err := fileImg.Write(pbImg.Img)

	if closeErr := fileImg.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		log.Printf("Image \"%s\" saved, %d bytes", fileName, size)
	}

	return err
}
