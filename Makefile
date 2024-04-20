include .env

.PHONY: build migrate migrate-down migrate-fix

install:
	go mod tidy

run-dev:
	air

run:
	go run main.go

build:
	go build -o bin/main main.go

run-build:
	./bin/main

up:
	docker compose up -d

down:
	docker compose down

migration:
	migrate create -seq -ext sql -dir db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate:
	migrate -path db/migrations -database "postgres://${DB_USER}:${DB_PWD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	migrate -path db/migrations -database "postgres://${DB_USER}:${DB_PWD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down $(filter-out $@,$(MAKECMDGOALS))

migrate-fix:
	migrate -path db/migrations -database "postgres://${DB_USER}:${DB_PWD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" force $(filter-out $@,$(MAKECMDGOALS))