build:
	docker compose -f deployments/docker-compose.yaml --env-file configs/config.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/config.env up
run: build up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/config.env down
lint:
	golangci-lint run ./...

.PHONY: build up run down lint

