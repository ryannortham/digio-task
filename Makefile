# Name of the binary file to be generated
BINARY_DIR=bin
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
EXT=$(if $(filter windows,$(GOOS)),.exe,)
BINARY_NAME=digio-task-$(GOOS)-$(GOARCH)$(EXT)

all: tidy fmt lint test build

build:
	go build -o $(BINARY_DIR)/$(BINARY_NAME) -v .

test: lint
	go test -v ./...

lint:
	golangci-lint run

clean:
	go clean

tidy:
	go mod tidy

fmt:
	go fmt ./...

run: build
	./$(BINARY_DIR)/$(BINARY_NAME) analyselog
