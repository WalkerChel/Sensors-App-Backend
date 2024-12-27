#!/usr/bin/env sh

# Export .env varibles to current terminal session
export $(cat ./.env | grep -v ^# | xargs) >/dev/null

echo "Successfully exported variables"
