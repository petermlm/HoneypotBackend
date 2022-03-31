#!/bin/bash

docker-compose exec postgres psql -h localhost -Uhoneypot_user -d honeypot_db
