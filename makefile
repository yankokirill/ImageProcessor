build:
	docker-compose build

run: build
	docker-compose up -d

tests: run
	docker-compose run --rm tests

stop:
	docker-compose down