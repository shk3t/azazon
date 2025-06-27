#!/usr/bin/env bash

if [[ $(basename $PWD) == "azazon" ]]; then
    cd base
fi

protoc \
    --proto_path=./api/proto \
    --go_out=paths=source_relative:./api/go \
    --go-grpc_out=paths=source_relative:./api/go \
    ./api/proto/*.proto
