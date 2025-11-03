#!/usr/bin/env bash

source $(dirname $0)/envs.sh

api_path="$(dirname $(dirname $0))/api"

for module in ${modules[@]}; do
    mkdir -p $api_path/${module}
    protoc \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --proto_path=$api_path/proto \
        --go_out=$api_path/${module} \
        --go-grpc_out=$api_path/${module} \
        $api_path/proto/${module}.proto
done
