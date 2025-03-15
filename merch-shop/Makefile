build:
	docker compose -f deployments/docker-compose.yaml --env-file configs/postgres.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/postgres.env up
run: build up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/postgres.env down
lint:
	golangci-lint run ./...

.PHONY: build up run down lint

