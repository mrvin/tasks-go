build:
	go build -o top-product -ldflags '-w -s'
lint:
	golangci-lint run ./...
docker-build:
	docker build -t top-product .
docker-run:
	docker run --rm --name top-product top-product
