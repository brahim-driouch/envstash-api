.PHONY: run test help

run:
	@go run ./cmd/api/main.go

test:
	@go test -v ./...

build:
	@go build -o bin/api ./cmd/api/main.go

help: