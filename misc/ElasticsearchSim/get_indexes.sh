#!/bin/bash

curl -iv --raw -XGET 'localhost:9250/_cat/indices?v'
