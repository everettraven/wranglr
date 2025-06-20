.PHONY: lint
lint:
	go tool golangci-lint run

build:
	go build -o synkr main.go

install:
	go install

.PHONY: test
test:
	@echo "TODO: Add testing :)"
