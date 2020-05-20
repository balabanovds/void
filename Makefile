CURRENT_USER = $(shell id -u):$(shell id -g)
export CURRENT_USER

.PHONY: run
run:
	go run ./cmd/void

pg-shell:
	docker exec -it void_posgres_1 bash

docker-up:
	docker-compose up

echo: 
	echo $(CURRENT_USER)

DEFAULT := run