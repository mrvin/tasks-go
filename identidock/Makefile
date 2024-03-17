build:
	docker-compose -f deployments/docker-compose.yaml build
up:
	docker-compose -f deployments/docker-compose.yaml up
down:
	docker-compose -f deployments/docker-compose.yaml down
lint:
	golangci-lint run ./...

.PHONY: build up down lint

