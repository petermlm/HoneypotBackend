# HoneypotBackend

This is the Go backend of an implementation of a honeypot. This honeypot will
open a few tcp ports, listen on them, and register any connection attempt to a
timeline. The timeline is kept in a Postgres database.

The Honeypot is implemented using a microservices architecture. The services
communicate via RabbitMQ.

# Honeypot's Services

## Listener

Opens ports, waits for connections, and in come cases, may pretend to be a
service to full the actor trying to access the compromised service.

## Processor

Takes in information from attempted connections and builds information about
them, such as the location of the connection.

## Storer

Takes in information from the processor and stores them in the database

## Webserver

A simple HTTP server that implements endpoints to be consumed by a frontend,
such as a web frontend that displays the data in a readable or interesting format.

# Run

## Docker-compose

Everything can be executed simply doing:

    docker-compose up

Or individual services may be started with:

    docker-compose up listener

## Start using go

In `src` there is a `Makefile`. Things can be run from it. This will execute go
directly without a docker container. For this, RabbitMQ and Influx can still be
started using docker-compose. The host for InfluxDB and RabbitMQ should be
changed to "localhost" in the config.

# Tests

In `src`, run `make test`. This will execute every unit test.

# Directories

## docker

All docker files and helper scripts for docker are placed here.

## misc

Contains some experiments and prototypes that go in the direction of making an
Elasticsearch simulator. The idea is to simulate a database (or any other
service) to full an attacker into thinking he is interacting with a real
exposed service. Like this, one can observe and study how an attacker works.
