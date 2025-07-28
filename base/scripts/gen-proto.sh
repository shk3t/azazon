#!/usr/bin/env bash

if [[ $(basename $PWD) == "azazon" ]]; then
    cd base
fi

modules=(auth)

for module in $modules; do
    mkdir -p api/${module}
    protoc \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --proto_path=./api/proto \
        --go_out=./api/auth \
        --go-grpc_out=./api/auth \
        ./api/proto/auth.proto
done
