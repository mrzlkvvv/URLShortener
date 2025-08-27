.SILENT:

COMPOSE_PATH = ./deploy/docker-compose.yml
COMPOSE_DEV_OVERRIDE_PATH = ./deploy/docker-compose.dev-override.yml

up-dev:
	docker compose -f $(COMPOSE_PATH) -f $(COMPOSE_DEV_OVERRIDE_PATH) build
	docker compose -f $(COMPOSE_PATH) -f $(COMPOSE_DEV_OVERRIDE_PATH) up
	docker compose -f $(COMPOSE_PATH) -f $(COMPOSE_DEV_OVERRIDE_PATH) down

up-prod:
	docker compose -f $(COMPOSE_PATH) build
	docker compose -f $(COMPOSE_PATH) up
	docker compose -f $(COMPOSE_PATH) down

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run ./{cmd,internal}/... --config=./config/.golangci.yml

test:
	go test -vet=off ./...
