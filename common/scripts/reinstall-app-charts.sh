#!/usr/bin/env bash

source $(dirname $0)/envs.sh

for app in ${apps[@]}; do
    app_dir="${app}/deployments/${app}"
    if [[ -d ${app_dir} ]]; then
        helm uninstall -n $namespace $app --ignore-not-found
        helm install -n $namespace $app $app_dir || exit 1
        kubectl rollout status -n $namespace deployment/${app}-deployment --timeout=10s || exit 1
    fi
done
