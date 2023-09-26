build:
	go build -o ../../bin/server-imgstorage -ldflags '-w -s'

run:
	go run -race main.go -config ../configs/imgstorage.yml

lint:
	golangci-lint run ./...

