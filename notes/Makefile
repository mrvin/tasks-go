build:
	docker compose -f deployments/docker-compose.yaml --env-file configs/notes.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/notes.env up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/notes.env down

.PHONY: build up down

