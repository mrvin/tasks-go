build:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env up
run: build up
down:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env down
lint:
	golangci-lint run ./...
unit-test:
	cd internal/httpserver/handlers/ && go test ./...

.PHONY: build up run down lint unit-test

