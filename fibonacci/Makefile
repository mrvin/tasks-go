build:
	docker-compose -f deployments/docker-compose.yml build
up:
	docker-compose -f deployments/docker-compose.yml up
down:
	docker-compose -f deployments/docker-compose.yml down
lint:
	golangci-lint run ./internal/...
codegen:
	go generate ./...

.PHONY: build up down lint codegen
