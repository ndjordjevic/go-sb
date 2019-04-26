#!/bin/bash

# Generate protobuf stubs. Run this from the project root
protoc api/order.proto --go_out=plugins=grpc:.
