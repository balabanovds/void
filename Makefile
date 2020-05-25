# use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)

include .env
DATABASE_URL = "host=$(PG_HOST) port=$(PG_PORT) user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DATABASE_TEST) sslmode=disable"

run:
	go run ./cmd/void

tests:
	DATABASE_URL=$(DATABASE_URL) go test  -v -race ./...

shell-pg:
	docker exec -it void_posgres_1 bash

migrate-test-reset: migrate-test-down migrate-test-up

migrate-test-up:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" up $(RUN_ARGS)

migrate-test-down:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE_TEST)?sslmode=disable" down $(RUN_ARGS)

migrate-up:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database \
		"postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=disable" down

docker-up:
	docker-compose up


DEFAULT := run