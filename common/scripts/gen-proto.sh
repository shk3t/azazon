#!/usr/bin/env bash

if [[ $(basename $PWD) == "azazon" ]]; then
    cd common
fi

modules=(common auth notification order payment stock)

for module in ${modules[@]}; do
    mkdir -p api/${module}
    protoc \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --proto_path=./api/proto \
        --go_out=./api/${module} \
        --go-grpc_out=./api/${module} \
        ./api/proto/${module}.proto
done
