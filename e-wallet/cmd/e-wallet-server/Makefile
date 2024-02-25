build:
	go build -o ../../bin/e-wallet-server -ldflags '-w -s'
test:
	go test ../../internal/httpserver/handlers/wallet/... -cover | grep "coverage:"

.PHONY: build test
