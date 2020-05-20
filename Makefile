include .env

run:
	go run ./cmd/void

pg-shell:
	docker exec -it void_posgres_1 bash

test-migrate-up:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" up

test-migrate-down:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" down

docker-up:
	docker-compose up


DEFAULT := run