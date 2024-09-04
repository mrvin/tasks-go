LDFLAGS := -w -s
build:
	go build -o ../../bin/notes-server -ldflags "$(LDFLAGS)"
lint:
	golangci-lint run ../../...

.PHONY: build lint
