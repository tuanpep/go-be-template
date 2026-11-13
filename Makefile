.PHONY: help run test test-race lint generate build compose-up compose-down migrate tidy docker-build

LOCAL_DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/go_be_template?sslmode=disable

GOBIN := $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)

help:
	@echo "Targets: run, test, test-race, lint, generate, build, migrate, compose-up"

compose-up:
	docker compose up -d --wait

compose-down:
	docker compose down

generate:
	go tool sqlc generate
	go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate

migrate:
	@bash -eu -c 'set -a; [ -f .env ] && . ./.env; set +a; export DATABASE_URL="$${DATABASE_URL:-$(LOCAL_DATABASE_URL)}"; go run ./cmd/migrate'

run: compose-up
	@bash -eu -c 'set -a; [ -f .env ] && . ./.env; set +a; export DATABASE_URL="$${DATABASE_URL:-$(LOCAL_DATABASE_URL)}"; go run ./cmd/migrate && go run ./cmd/api'

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	@which golangci-lint >/dev/null || (echo "install golangci-lint v2.9.0: https://golangci-lint.run/welcome/install/" && exit 1)
	golangci-lint run ./...

docker-build:
	docker build -t go-be-template:local .
