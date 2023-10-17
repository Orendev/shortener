#!/bin/bash
# run it from makefile
protoc \
  --proto_path=./proto \
  --go_out=. \
  --go_opt=paths=import \
  --go-grpc_out=. \
  --go-grpc_opt=paths=import ./proto/shortener.proto