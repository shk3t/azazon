#!/usr/bin/env bash

source $(dirname $0)/envs.sh

for app in ${apps[@]}; do
    db_dir="${app}/deployments/${app}-db"
    if [[ -d ${db_dir} ]]; then
        helm uninstall -n $namespace ${app}-db --ignore-not-found
        helm install -n $namespace ${app}-db $db_dir || exit 1
        kubectl rollout status -n $namespace statefulset/${app}-database-statefulset --timeout=10s || exit 1
    fi
done
