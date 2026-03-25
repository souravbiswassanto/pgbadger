BINARY_NAME=pgbadger

.PHONY: all build run tidy test docker

all: build

tidy:
	go mod tidy

build:
	go build -o $(BINARY_NAME) ./...

run: build
	./$(BINARY_NAME) server

test:
	go test ./...

docker:
	docker build -t $(BINARY_NAME):latest .
