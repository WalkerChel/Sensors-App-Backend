#!/usr/bin/env bash

MIGRATIONS_PATH=$1
TYPE=$2

case "$TYPE" in
    'init')
        echo "Versioning migrations++"
        migrate create -ext sql -dir "${MIGRATIONS_PATH}" -seq init_schema
    ;;
    'up')
            echo "Migrating up..."
            migrate -path "${MIGRATIONS_PATH}" -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" -verbose up
    ;;
    'down')
            echo "Migrating down..."
            migrate -path "${MIGRATIONS_PATH}" -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" -verbose down
esac
