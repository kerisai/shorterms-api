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