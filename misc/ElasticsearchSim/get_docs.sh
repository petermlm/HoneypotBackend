#!/bin/bash

curl -iv --raw -XGET 'localhost:9250/users/_search?pretty=true&q=*:*'
