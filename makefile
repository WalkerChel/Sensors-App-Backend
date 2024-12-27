include ./.env

MIGRATIONS_PATH=./migrations
INIT_WAY=init
WAY_UP=up
WAY_DOWN=down

.PHONY: init-vars
init-vars:
	./scripts/env.sh  

.PHONY:postgres-up
postgres-up:
	docker run --name=${POSTGRES_DB} \
	-e POSTGRES_USER=${POSTGRES_USER} \
	-e POSTGRES_DB=${POSTGRES_DB} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-e POSTGRES_SSL_MODE=${POSTGRES_SSL_MODE} \
	-p=${POSTGRES_PORT}:${POSTGRES_PORT} \
	-d postgres

.PHONY:postgres-down
postgres-down:
	docker stop ${POSTGRES_DB} && docker rm ${POSTGRES_DB};

.PHONY: migrate-init
migrate-init:
	@make init-vars && ./scripts/migrations.sh ${MIGRATIONS_PATH} ${INIT_WAY}

.PHONY: migrate-up
migrate-up:
	@make init-vars && ./scripts/migrations.sh ${MIGRATIONS_PATH} ${WAY_UP}

.PHONY: migrate-down
migrate-down:
	@make init-vars && ./scripts/migrations.sh ${MIGRATIONS_PATH} ${WAY_DOWN}

.PHONY: help
help:
	@echo "												"
	@echo "	available commands:							"
	@echo "												"
	@echo "	postgres-up	->		creates postgres db docker container"
	@echo "	postgres-down	->		stops and removes postgres db docker container"
	@echo "	migrate-init	->		creates a new init blank files"
	@echo "	migrate-up	->		applies up migrations to db"
	@echo "	migrate-down	->		applies down migrations to db"
