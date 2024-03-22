build:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env up
down:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env down
lint:
	golangci-lint run ./...

.PHONY: build up down lint

