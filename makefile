include .env
export

.PHONY: postgres-up
postgres-up:
	docker compose up -d postgres_db

.PHONY: postgres-down
postgres-down:
	docker compose down postgres_db

.PHONY: redis-up
redis-up:
	docker compose up -d redis_db

.PHONY: redis-down
redis-down:
	docker compose down redis_db

.PHONY: migrate-up
migrate-up:
	@echo "Migrating up..."
	docker compose run --rm migrator

.PHONY: migrate-down
migrate-down:
	@echo "Migrating down..."
	docker compose run --rm migrator -- --path migrations --database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres_db:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" down

.PHONY: app-up
app-up:
	docker compose up --build -d app

.PHONY: app-down
app-down:
	docker compose down app

.PHONY: infrastructure-up
infrastructure-up:
	docker compose up -d postgres_db redis_db

.PHONY: infrastructure-stop
infrastructure-stop:
	docker compose down postgres_db redis_db

.PHONY: compose-up
compose-up:
	docker compose up --build -d

.PHONY: compose-down
compose-down:
	docker compose down -v

.PHONY: help
help:
	@echo "                                                                               "
	@echo " available commands:                                                           "
	@echo "                                                                               "
	@echo " compose-up          ->    builds and starts all docker-compose services"
	@echo " compose-down        ->    stops and removes all services and volumes"
	@echo " infrastructure-up   ->    starts postgres and redis databases"
	@echo " infrastructure-stop ->    stops and removes all databases / services"
	@echo " app-up              ->    builds and runs the app container"
	@echo " app-down            ->    stops and removes the app container"
	@echo " postgres-up         ->    creates postgres db container via compose"
	@echo " postgres-down       ->    stops and removes postgres db container"
	@echo " redis-up            ->    creates redis db container via compose"
	@echo " redis-down          ->    stops and removes redis db container"
	@echo " migrate-up          ->    runs up migrations using the compose migrator service"
	@echo " migrate-down        ->    runs down migrations against the database"
	@echo "                                                                               "
