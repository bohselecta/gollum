.PHONY: run build test tidy

run:
	go run ./cmd/altiserve

build:
	go build -o bin/altiserve ./cmd/altiserve

test:
	go test ./...

tidy:
	go mod tidy
