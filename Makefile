stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "no containers running..." ;\
	fi

run-db:
	docker run --name mongodb -d -p 27017:27017 mongo:latest

build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :7071

test:
	@go test -v ./...