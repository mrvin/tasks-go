build:
	go build -o go-calc -ldflags '-w -s'
test:
	go test ./...
cover:
	go test ./... -cover
lint:
	golangci-lint run .
