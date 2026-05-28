.PHONY: list install run run-web test test-e2e test-api build build-web format vet tidy lint upgrade compose-start compose-stop compose-restart
help:
	@grep -E '^[[:alpha:]].*:' Makefile | cat # Get all targets in this file, without color-coding the matching letters

install:
	go install ./...

# Filter out known targets so extra args are passed through to the binary
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

run:
	go run ./cmd/webscan/... $(RUN_ARGS)

run-web:
	go run ./cmd/webscan-web/... $(RUN_ARGS)

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

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	go vet ./...
	go build ./...
	npm run lint:md

upgrade:
	go get -u ./...
	go mod tidy
