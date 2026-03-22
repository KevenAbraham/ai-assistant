.PHONY: build build-api build-daemon run-api run-daemon migrate-up migrate-down tidy vet test docker-up docker-down

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
