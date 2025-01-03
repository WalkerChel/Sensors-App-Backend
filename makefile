include .env
export

MIGRATIONS_PATH=./migrations

.PHONY:postgres-up
postgres-up:
	docker run --name=${POSTGRES_DB} \
	-e POSTGRES_USER=${POSTGRES_USER} \
	-e POSTGRES_DB=${POSTGRES_DB} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-e POSTGRES_SSL_MODE=${POSTGRES_SSL_MODE} \
	-p=${POSTGRES_PORT}:${POSTGRES_PORT} \
	-d postgres:17.2

.PHONY:postgres-down
postgres-down:
	docker stop ${POSTGRES_DB} && docker rm ${POSTGRES_DB};

.PHONY:redis-up
redis-up:
	docker run --name=${REDIS_CONTAINER_NAME} \
	-e REDIS_USER=${REDIS_USER} \
	-e REDIS_USER_PASSWORD=${REDIS_USER_PASSWORD} \
	-p=${REDIS_PORT}:${REDIS_PORT} \
	-d redis:7.4.1

.PHONY:redis-down
redis-down:
	docker stop ${REDIS_CONTAINER_NAME} && docker rm ${REDIS_CONTAINER_NAME};

.PHONY: migrate-up
migrate-up:
	echo "Migrating up..."
	migrate -path ${MIGRATIONS_PATH} -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}' -verbose up


.PHONY: migrate-down
migrate-down:
	echo "Migrating down..."
	migrate -path ${MIGRATIONS_PATH} -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}' -verbose down

.PHONY: infrastructure-stop
infrastructure-stop:
	@make postgres-down \
	&& make redis-down

.PHONY: infrastructure-up
infrastructure-up:
	@make postgres-up \
	&& make redis-up

.PHONY: help
help:
	@echo "												"
	@echo "	available commands:							"
	@echo "												"
	@echo "	infrastructure-up	->		start all databases"
	@echo "	infrastructure-stop	->		remove all databases"
	@echo "	postgres-up	->		creates postgres db docker container"
	@echo "	postgres-down	->		stops and removes postgres db docker container"
	@echo "	redis-up	->		creates redis db docker container"
	@echo "	redis-down	->		stops and removes redis db docker container"
	@echo "	migrate-up	->		applies up migrations to db"
	@echo "	migrate-down	->		applies down migrations to db"
