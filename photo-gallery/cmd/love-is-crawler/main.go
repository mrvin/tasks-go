package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const url = "https://qr.loveis.site/iid%d/lang3.png"

const idOffset = 7

func main() {
	for i := 8; i < 84; i++ {
		imageName := fmt.Sprintf("../../testdata/love_is_%d.png", i-idOffset)
		sizeImage, err := getImage(fmt.Sprintf(url, i), imageName)
		if err != nil {
			log.Printf("Failure get image: %v\n", err)
		}
		log.Printf("Downloaded a image %s with size %d\n", imageName, sizeImage)
	}
}

func getImage(url, filmName string) (int64, error) {
	resp, err := http.Get(url) //nolint:gosec,noctx
	if err != nil {
		return 0, fmt.Errorf("get http request: %w", err)
	}
	defer resp.Body.Close()

	filePoster, err := os.Create(filmName)
	if err != nil {
		return 0, fmt.Errorf("create image file: %w", err)
	}
	defer filePoster.Close()

	size, err := io.Copy(filePoster, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("copy image from response body to file: %w", err)
	}

	return size, nil
}
