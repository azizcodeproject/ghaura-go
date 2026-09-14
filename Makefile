.PHONY: up down infra api logs test-api

up:
	docker compose up --build

infra:
	docker compose up -d postgres redis redpanda

down:
	docker compose down

api:
	cd apps/api && go run ./cmd/server

logs:
	docker compose logs -f

test-api:
	cd apps/api && go test ./...
