#!/bin/bash
set -e
cd /app
protoc \
  -I /app \
  -I /usr/local/include/googleapis \
  --go_out=/app --go_opt=module=golang_restapi/contracts \
  --go-grpc_out=/app --go-grpc_opt=module=golang_restapi/contracts \
  $(find account auth common gateway pagination transaction -name '*.proto')