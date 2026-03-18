#!/usr/bin/env bash
sleep 3
goose postgres "host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB" up