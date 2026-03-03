#!/bin/bash

set -e

echo "Running alembic database migrations"
alembic upgrade head

IS_DEBUG=${IS_DEBUG:-True}

if [ "$IS_DEBUG" = True ]; then
    echo "Running Identity-Service in dev mode"
    uvicorn src.main:app --host 0.0.0.0 --port 80 --reload --workers 1
else
    echo "Running Identity-Service in prod mode"
    uvicorn src.main:app --host 0.0.0.0 --port 80 --workers 4
fi