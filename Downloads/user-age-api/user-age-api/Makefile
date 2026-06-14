.PHONY: run build test sqlc migrate-up migrate-down docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path db/migrations -database "$$DATABASE_URL" down

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
