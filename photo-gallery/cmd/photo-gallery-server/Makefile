build:
	go build -o ../../bin/photo-gallery-server -ldflags '-w -s'
lint:
	cd ../../ && golangci-lint run ./...

.PHONY: build lint
