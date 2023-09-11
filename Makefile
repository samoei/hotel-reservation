build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :7070

test:
	@go test -v ./...