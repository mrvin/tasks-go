build:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env down

.PHONY: build up down

