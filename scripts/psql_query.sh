#!/bin/bash

docker-compose exec postgres psql -h localhost -Uhoneypot_user -d honeypot_db -c 'select time, port, ip, country_code, client_port from conn_attemps'
