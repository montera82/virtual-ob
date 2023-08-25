.PHONY: all deps test docker

all: deps test docker

deps:
	@go mod download

test:
	@docker run --rm -v "$(CURDIR)":/app -w /app golang:1.20 go test -v -race -cover ./...

docker:
	@docker-compose up --build -d && docker-compose logs -f virtual-orb
