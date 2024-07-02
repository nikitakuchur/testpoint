.DEFAULT_GOAL := build

.PHONY: build clean
build:
	go fmt ./...
	go build -o ./bin/testpoint
clean:
	go clean
	rm -rf ./bin
