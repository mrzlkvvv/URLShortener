.SILENT:

run: build
	./bin/URLShortener

build:
	go build -o ./bin/URLShortener ./cmd/main.go

update:
	go get -u ./...
	go mod tidy
