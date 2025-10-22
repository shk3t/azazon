#!/usr/bin/env bash

set -a
source .env
set +a

apps=(common auth notification order payment stock)

for app in ${apps[@]}; do
    APP_DIR="${app}/deployments/${app}"
    if [[ -d ${APP_DIR} ]]; then
        envsubst < ${APP_DIR}/values.yaml.tpl > ${APP_DIR}/values.yaml
    fi
    if [[ -d ${APP_DIR}-db ]]; then
        envsubst < ${APP_DIR}-db/values.yaml.tpl > ${APP_DIR}-db/values.yaml
    fi
done
