.PHONY: lint
lint:
	go tool golangci-lint run

build:
	go build -o ./bin/wranglr main.go

install:
	go install
