#!/usr/bin/env bash

#Start elasticsearch from ~ dir
elasticsearch-7.0.0/bin/elasticsearch

# Test elasticsearch installation
curl http://localhost:9200/

# Check cluster health
curl -X GET "localhost:9200/_cat/health?v"

# List all indexes
curl -X GET "localhost:9200/_cat/indices?v"

# DELETE a index
curl -X DELETE "localhost:9200/orders?pretty"
