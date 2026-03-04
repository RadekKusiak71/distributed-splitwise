#!/bin/sh

set -e

echo "Running database migrations"
goose up

IS_DEBUG=${IS_DEBUG:-True}

if [ "$IS_DEBUG" = True ]; then
    echo "Running API-Gateway-Service in dev mode"
    make run-dev
else
    echo "Running API-Gateway-Service in prod mode"
    make run-prod
fi