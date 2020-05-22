include .env

run:
	go run ./cmd/void

tests:
	go test -v -race ./...

shell-pg:
	docker exec -it void_posgres_1 bash

migrate-test-reset: migrate-test-down migrate-test-up

migrate-test-up:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" up

migrate-test-down:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" down

migrate-up:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=disable" down

docker-up:
	docker-compose up


DEFAULT := run