.PHONY: list install run run-web test test-e2e test-api build build-web format lint upgrade compose-start compose-stop compose-restart
list:
	@grep -E '^[[:alpha:]].*:' Makefile | cat # Get all targets in this file, without color-coding the matching letters

install:
	go get ./...

run:
	go run ./cmd/webscan/...

run-web:
	go run ./cmd/webscan-web/...

test:
	go test -v ./...

test-e2e:
	cd tests && npm run test:e2e

test-api:
	cd tests && npm run test:api

compose-start:
	docker compose up -d --build

compose-stop:
	docker compose down

compose-restart:
	docker compose down
	docker compose up -d --build

build:
	go build -o webscan ./cmd/webscan/

build-web:
	go build -o webscan-web ./cmd/webscan-web/

format:
	go fmt ./...

lint:
	go vet ./...
	go build ./...
	npm run lint:md

upgrade:
	go get -u ./...
	go mod tidy
