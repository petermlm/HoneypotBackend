#!/bin/bash

echo localhost:8100/totalConsumptions
curl localhost:8100/totalConsumptions
echo

echo "----------------------------------------"
echo localhost:8100/map
curl localhost:8100/map
echo

echo "----------------------------------------"
echo localhost:8100/connAttemps
curl localhost:8100/connAttemps
echo

echo "----------------------------------------"
echo localhost:8100/topConsumers
curl localhost:8100/topConsumers
echo

echo "----------------------------------------"
echo localhost:8100/topFlavours
curl localhost:8100/topFlavours
echo

echo "----------------------------------------"
echo localhost:8100/getBytes/postgresql
curl localhost:8100/getBytes/postgresql
echo localhost:8100/getBytes/mysql
curl localhost:8100/getBytes/mysql
