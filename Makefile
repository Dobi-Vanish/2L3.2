.PHONY: build run test migrate

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

migrate:
	go run ./cmd/migrations/migrations.go

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down