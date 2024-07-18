.DEFAULT_GOAL := build

.PHONY: build clean test
build:
	go fmt ./...
	go build -o ./bin/testpoint -ldflags="-X main.version=$(VERSION)" ./cmd/testpoint
clean:
	go clean
	rm -rf ./bin
test:
	go test -v ./...
