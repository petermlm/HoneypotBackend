#!/bin/bash

echo "Tearing down"
docker-compose down

echo "Building and starting"
docker-compose up -d --build

echo "Waiting to make migrations"
sleep 10;

echo "Migrating"
docker-compose exec storer go run cmd/migrate/migrate.go
