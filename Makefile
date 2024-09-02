.PHONY: build run test down logs ps

build:
	docker-compose build

run:
	docker-compose up -d

test:
	docker-compose run --rm app go test -v ./...


down:
	docker-compose down


logs:
	docker-compose logs -f


ps:
	docker-compose ps


up: build run


restart: down up