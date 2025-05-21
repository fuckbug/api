init: docker-down-clear docker-pull docker-build docker-up

up: docker-up
down: docker-down
restart: down up

docker-pull:
	docker compose -f ./deployments/development/docker-compose.yml pull

docker-build:
	docker compose -f ./deployments/development/docker-compose.yml build --pull

docker-up:
	docker compose -f ./deployments/development/docker-compose.yml up -d

docker-down:
	docker compose -f ./deployments/development/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker compose -f ./deployments/development/docker-compose.yml down -v --remove-orphans


migrations-new:
	migrate create -ext sql -dir ./internal/storage/sql/migrations -seq init


test:
	go test -race -count 100 ./internal/...


lint:
	docker compose -f ./deployments/development/docker-compose.yml run --rm golangci-lint

clean:
	go mod tidy

.PHONY: init up down restart migrations-new test lint clean