build:
	docker compose -f deployments/docker-compose.yaml build
up:
	docker compose -f deployments/docker-compose.yaml up
run: build up
down:
	docker compose -f deployments/docker-compose.yaml down
lint:
	golangci-lint run ./...

.PHONY: build up run down lint

