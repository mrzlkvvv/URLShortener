.SILENT:

BINARY_NAME = URLShortener
BINARY_PATH = ./bin/$(BINARY_NAME)

run: build
	DOTENV_PATH="./config/.env" $(BINARY_PATH)

build:
	go build -o $(BINARY_PATH) ./cmd/main.go

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run ./{cmd,internal}/... --config=./config/.golangci.yml

test:
	go test -vet=off ./...
