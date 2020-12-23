#!/bin/bash

# This script, or a variant thereof, should only be used for experiments in
# development purposes.

echo -ne "This is some random text" | netcat localhost 5432
