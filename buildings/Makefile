build:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env up
down:
	docker compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env down
gen_docs:
	swag init -g internal/httpserver/server.go

.PHONY: build up down gen_docs
