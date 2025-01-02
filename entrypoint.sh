#!/bin/bash

set -e

# Export .env varibles to current terminal session
export $(cat ./.env | grep -v ^# | xargs) >/dev/null

if [ -z "${MIGRATIONS_PATH_DOCKER}" ]; then
    echo "Missed migration path variable. Exiting ..."
    exit 1
fi

  echo "Starting migrations"
  migrate -path ${MIGRATIONS_PATH_DOCKER} -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_DB}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" -verbose up

echo "Starting service"
exec ./SensorsApp
