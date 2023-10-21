# Name of the binary file to be generated
BINARY_NAME=digio-task.exe

# Lint parameters
LINTCMD=golangci-lint
LINTCMDARGS=run

all: tidy fmt lint test build

build:
	go build -o $(BINARY_NAME) -v .

test: lint
	go test -v ./...

lint:
	$(LINTCMD) $(LINTCMDARGS)

clean:
	go clean

tidy:
	go mod tidy

fmt:
	go fmt ./...

run: build
	./$(BINARY_NAME) analyse
