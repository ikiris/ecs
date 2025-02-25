all: build test benchmark

build:
	go build -v ./...

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

benchmark:
	go test -v -bench=. ./...

coverage:
	go tool cover -html=coverage.out
