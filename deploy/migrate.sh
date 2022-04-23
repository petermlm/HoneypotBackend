#!/bin/bash

docker-compose exec storer go cmd/migrate/migrate.go
