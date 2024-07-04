VERSION := 1.0.0

.DEFAULT_GOAL := build

.PHONY: build clean
build:
	go fmt ./...
	go build -o ./bin/testpoint -ldflags="-X main.version=${VERSION}" ./cmd/testpoint
clean:
	go clean
	rm -rf ./bin
