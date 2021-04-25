#!/bin/bash

docker-compose \
    -p honeypot \
    -f docker-compose-go.yml \
    -f docker-compose-influxdb.yml \
    -f docker-compose-influxdb-tooling.yml \
    -f docker-compose-rabbitmq.yml \
    down
