build:
	go build -o ../../bin/e-wallet-server -ldflags '-w -s'
test:
	go test ../../internal/httpserver/handlers/wallet/... -cover | grep -v "mocks"

.PHONY: build test
