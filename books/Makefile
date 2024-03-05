test:
	go test -race ./internal/storage/... -cover
	cd services/calendar && make test

build:
	docker-compose -f deployments/docker-compose.yaml build

up:
	docker-compose -f deployments/docker-compose.yaml up

down:
	docker-compose -f deployments/docker-compose.yaml down

lint:
	golangci-lint run ./internal/...

codegen:
	go generate ./...

.PHONY: test build up down lint codegen

