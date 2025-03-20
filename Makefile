APP_NAME = mortgage-calc-service
DOCKER_IMAGE = $(APP_NAME)
DOCKER_CONTAINER = mortgage-calc

test:
	go test ./... -v

lint:
	golangci-lint run ./...

build:
	docker build -t $(DOCKER_IMAGE) .

run:
	docker run --rm -d --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE)

stop:
	docker stop $(DOCKER_CONTAINER) || true

clean:
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

dev: build run

all: lint test build run
