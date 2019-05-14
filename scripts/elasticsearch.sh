#!/usr/bin/env bash

#Start elasticsearch from ~ dir
elasticsearch-7.0.0/bin/elasticsearch

# Test elasticsearch installation
curl http://localhost:9200/

# Check cluster health
curl -X GET "localhost:9200/_cat/health?v"

# Nodes information
curl 'http://localhost:9200/_nodes/http?pretty'

# List all indexes
curl -X GET "localhost:9200/_cat/indices?v"

# DELETE a index
curl -X DELETE "localhost:9200/order?pretty"

# Search a term across all fields in an index
curl -XPOST -H 'Content-Type: application/json' 'localhost:9200/order/_search?pretty' -d '{
  "query": { "query_string": { "query": "ndjordjevic" } }
}'
