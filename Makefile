APP_NAME=discord-presence-bot
IMAGE=$(APP_NAME):latest

.PHONY: help tidy build run stop logs docker-build docker-run docker-stop compose up down restart

help:
	@echo "Targets:"
	@echo "  tidy           - go mod tidy"
	@echo "  build          - go build local binary"
	@echo "  run            - run locally with .env"
	@echo "  docker-build   - build docker image"
	@echo "  docker-run     - run container (needs .env)"
	@echo "  docker-stop    - stop container"
	@echo "  logs           - tail docker logs"
	@echo "  up             - docker compose up -d"
	@echo "  down           - docker compose down"
	@echo "  restart        - docker compose restart"

tidy:
	go mod tidy

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/bot .

run:
	@export $$(grep -v '^#' .env | xargs) && ./bin/bot

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	@if [ ! -f .env ]; then echo ".env missing (copy .env.example)"; exit 1; fi
	docker run --name $(APP_NAME) --rm -d --env-file .env $(IMAGE)

docker-stop:
	-@docker stop $(APP_NAME) >/dev/null 2>&1 || true

logs:
	docker logs -f $(APP_NAME)

up:
	docker compose up -d --build

down:
	docker compose down

restart:
	docker compose restart
