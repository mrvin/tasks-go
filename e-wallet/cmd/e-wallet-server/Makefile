LDFLAGS := -w -s
build:
	go build -o ../../bin/e-wallet-server -ldflags "$(LDFLAGS)"
test:
	go test ../../internal/httpserver/handlers/wallet/... -cover | grep -v "mocks"

.PHONY: build test
