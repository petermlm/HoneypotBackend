#!/bin/bash

docker-compose \
    -p honeypot \
    -f docker-compose-influxdb.yml \
    -f docker-compose-influxdb-tooling.yml \
    down
