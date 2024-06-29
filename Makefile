.DEFAULT_GOAL := build

.PHONY: build clean
build:
	go fmt ./...
	go build -o ./bin/restcompare
clean:
	go clean
	rm -rf ./bin
