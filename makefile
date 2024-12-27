include ./.env

MIGRATIONS_PATH=./migrations
INIT_WAY=init
WAY_UP=up
WAY_DOWN=down

.PHONY: init-vars
init-vars:
	./scripts/env.sh  

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
	@echo "	available commands:							"
	@echo "												"
	@echo "	migrate-init	->		creates a new init blank files"
	@echo "	migrate-up	->		applies up migrations to db"
	@echo "	migrate-down	->		applies down migrations to db"
