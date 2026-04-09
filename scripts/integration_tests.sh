#!/bin/bash

set -e

export PGTEST_EXTERNAL=true
export PGTEST_IMAGE=postgres:17.4
export PGTEST_CONTAINER_NAME=test-db-media-equipment-tracker

export PGTEST_MIGRATION_DIRECTORY=./deployments/migrations
export PGTEST_WORKING_DIRECTORY=$(pwd)
export PGTEST_USERNAME=postgres
export PGTEST_PASSWORD=postgres
export PGTEST_DATABASE=default
export PGTEST_PORT=35432

docker compose -f ./deployments/docker-compose.test.yaml up -d
CONTAINER_ID=$(docker compose -f ./deployments/docker-compose.test.yaml ps -q ${PGTEST_CONTAINER_NAME})
if [ -z "$CONTAINER_ID" ]; then
    echo "Container ${PGTEST_CONTAINER_NAME} not found"
    exit 1
fi

TIMEOUT=60
ELAPSED=0

while [ $ELAPSED -lt $TIMEOUT ]; do
    STATUS=$(docker inspect --format='{{.State.Health.Status}}' $CONTAINER_ID 2>/dev/null)
    if [ "$STATUS" = "healthy" ]; then
        echo "Container is healthy!"
        break
    elif [ "$STATUS" = "unhealthy" ]; then
        echo "Container is unhealthy!"
    fi
    sleep 2
    ELAPSED=$((ELAPSED + 2))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
    echo "Container did not become healthy within ${TIMEOUT} seconds"
    exit 1
fi

go test -v -tags=integration ./...
