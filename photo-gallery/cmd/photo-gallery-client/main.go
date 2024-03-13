package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var port int
	var host string
	var pathToFile string

	flag.IntVar(&port, "port", 8088, "port")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.StringVar(&pathToFile, "name", "", "path to file")
	flag.Parse()

	if pathToFile == "" {
		log.Print("Path to file is empty")
		return
	}
	filePhoto, err := os.Open(pathToFile)
	if err != nil {
		log.Print("Open photo: %v", err)
		return
	}
	defer filePhoto.Close()

	_, fileName := filepath.Split(pathToFile)

	_, err = http.Post(fmt.Sprintf("http://%s:%d/api/v1/photo?name=%s", host, port, fileName), "image/jpeg", filePhoto)
	if err != nil {
		log.Print("HTTP post: %v", err)
		return
	}

}
