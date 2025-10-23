#!/usr/bin/env bash

set -a
source .env
set +a

apps=(auth notification order payment stock)

for app in ${apps[@]}; do
    app_dir="${app}/deployments/${app}"
    if [[ -d ${app_dir} ]]; then
        envsubst < ${app_dir}/values.yaml.tpl > ${app_dir}/values.yaml
    fi
    if [[ -d ${app_dir}-db ]]; then
        envsubst < ${app_dir}-db/values.yaml.tpl > ${app_dir}-db/values.yaml
    fi
done
