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

generate-docs:
	swag init -g cmd/fuckbug/main.go

clean:
	go clean -modcache && go mod tidy

docker-pull:
	docker compose -f ./deployments/development/docker-compose.yml pull

docker-build:
	docker compose -f ./deployments/development/docker-compose.yml build --pull

dockerhub-build:
	docker build --platform=linux/amd64 -f ./build/fuckbug/Dockerfile -t fuckbugio/api:1.0.0 .

dockerhub-push:
	docker push fuckbugio/api:1.0.0

dockerhub-deploy: dockerhub-build dockerhub-push

docker-up:
	docker compose -f ./deployments/development/docker-compose.yml up -d

docker-down:
	docker compose -f ./deployments/development/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker compose -f ./deployments/development/docker-compose.yml down -v --remove-orphans

docker-network:
	docker network create fuckbug_network || true

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/fuckbug

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version


migrations-new:
	migrate create -ext sql -dir ./internal/storage/sql/migrations -seq init


test:
	go test -race -count 100 ./internal/...


remove-lint-deps:
	rm $(which golangci-lint)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.2

lint: install-lint-deps
	golangci-lint run

lint-fix: install-lint-deps
	golangci-lint run --fix


.PHONY: build run build-img run-img version test lint