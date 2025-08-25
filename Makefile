.SILENT:

BINARY_NAME = URLShortener
BINARY_PATH = ./bin/$(BINARY_NAME)

run: build
	CONFIG_PATH="./config/config.yml" $(BINARY_PATH)

build:
	go build -o $(BINARY_PATH) ./cmd/main.go

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run ./... --config=./config/.golangci.yml

test:
	go test -vet=off ./...
