#!/bin/bash

docker-compose exec storer go run cmd/migrate/migrate.go
