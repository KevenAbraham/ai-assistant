.PHONY: build build-api build-daemon run-api run-daemon migrate-up migrate-down tidy vet test docker-up docker-down

export GOTOOLCHAIN := path

_GO_1_25_BIN := /usr/local/go1.25.1/bin
ifneq ($(wildcard $(_GO_1_25_BIN)/go),)
export PATH := $(_GO_1_25_BIN):$(PATH)
export GOROOT := /usr/local/go1.25.1
endif

build: build-api build-daemon

build-api:
	go build -o bin/api ./cmd/api

build-daemon:
	CGO_ENABLED=1 go build -o bin/daemon ./cmd/daemon

run-api:
	bash -c 'set -a && source .env 2>/dev/null; set +a && go run ./cmd/api'

run-daemon:
	bash -c 'set -a && source .env 2>/dev/null; set +a && CGO_ENABLED=1 go run ./cmd/daemon'

migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./...

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down
