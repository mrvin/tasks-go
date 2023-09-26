package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/mrvin/tasks-go/006-imgstorage/internal/imgstorageapi"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const chunkSize = 64 * 1024 // 64 KiB
const timeoutContext = 20 * time.Second

type Client struct {
	client imgstorageapi.ImgStorageClient
}

var ctx = context.Background()

func main() {
	var port int
	var host string
	var fileName string
	var upload, download, list bool

	flag.IntVar(&port, "port", 55555, "port")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.StringVar(&fileName, "name", "", "file name")
	flag.BoolVar(&upload, "upload", false, "upload image")
	flag.BoolVar(&download, "download", false, "download image")
	flag.BoolVar(&list, "list", false, "image list")
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

	switch {
	case upload:
		if err := c.uploadImg(fileName); err != nil {
			log.Printf("Upload image: %v", err)
		}
	case download:
		if err := c.downloadImg(fileName); err != nil {
			log.Printf("Download image: %v", err)
		}
	case list:
		if err := c.getListImg(); err != nil {
			log.Printf("Image list: %v", err)
		}
	}
}

func (c *Client) uploadImg(pathToFile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	stream, err := c.client.UploadImage(ctx)
	if err != nil {
		return fmt.Errorf("get stream: %w", err)
	}

	fileImage, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("open image file: %w", err)
	}
	defer fileImage.Close()

	fileName := filepath.Base(pathToFile)
	req := &imgstorageapi.UploadImageRequest{
		Data: &imgstorageapi.UploadImageRequest_Name{ //nolint: nosnakecase
			Name: fileName,
		},
	}

	if err := stream.Send(req); err != nil {
		return fmt.Errorf("send image name: %w", err)
	}

	reader := bufio.NewReader(fileImage)
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read chunk to buffer: %w", err)
		}

		req := &imgstorageapi.UploadImageRequest{
			Data: &imgstorageapi.UploadImageRequest_ChunkData{ //nolint: nosnakecase
				ChunkData: buffer[:n],
			},
		}

		if err = stream.Send(req); err != nil {
			return fmt.Errorf("send chunk: %w", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close stream and receive response: %w", err)
	}

	log.Printf("Upload image \"%s\", %d byte", fileName, res.GetSize() /*, res.GetId()*/)

	return nil
}

func (c *Client) downloadImg(fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	req := &imgstorageapi.DownloadImageRequest{Name: fileName}

	stream, err := c.client.DownloadImage(ctx, req)
	if err != nil {
		return fmt.Errorf("get stream and send request: %w", err)
	}

	fileImage, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	defer fileImage.Close()

	imageSize := 0
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("receive chunk data: %w", err)
		}

		chunk := res.GetChunkData()

		size, err := fileImage.Write(chunk)
		if err != nil {
			return fmt.Errorf("write chunk data: %w", err)
		}

		imageSize += size
	}

	log.Printf("Image \"%s\" saved, %d bytes", fileName, imageSize)

	return nil
}

func (c *Client) getListImg() error {
	req := &emptypb.Empty{}

	pbListImg, err := c.client.GetListImage(ctx, req)
	if err != nil {
		return fmt.Errorf("get list image: %w", err)
	}

	const format = "%s\t%s\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "File name", "Modified date")
	fmt.Fprintf(tw, format, "---------", "-------------")

	for _, infImg := range pbListImg.ImageInfo {
		fmt.Fprintf(tw, format, infImg.Name, infImg.ModifiedAt.AsTime().Format("2 Jan 2006 15:04"))
	}

	tw.Flush()

	return nil
}
