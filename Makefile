BIN := "./bin/fuckbug"
DOCKER_IMG="fuckbug:fuckbug"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

init: docker-down-clear \
	docker-network \
	docker-pull docker-build docker-up

up: docker-network docker-up
down: docker-down
restart: down up

clean:
	go clean -modcache && go mod tidy

docker-pull:
	docker compose -f ./deployments/development/docker-compose.yml pull

docker-build:
	docker compose -f ./deployments/development/docker-compose.yml build --pull

build:
	docker build --platform=linux/amd64 -f ./build/fuckbug/Dockerfile -t fuckbugio/api:1.0.0 .

push:
	docker push fuckbugio/api:1.0.0

deploy: build push

docker-up:
	docker compose -f ./deployments/development/docker-compose.yml up -d

docker-down:
	docker compose -f ./deployments/development/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker compose -f ./deployments/development/docker-compose.yml down -v --remove-orphans

docker-network:
	docker network create fuckbug_network || true


migrations-new:
	migrate create -ext sql -dir ./internal/storage/sql/migrations -seq init


test:
	go test -race -count 100 ./internal/...


lint:
	docker compose -f ./deployments/development/docker-compose.yml run --rm golangci-lint


.PHONY: test lint