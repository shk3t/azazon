#!/usr/bin/env bash

source $(dirname $0)/envs.sh

eval $(minikube -p minikube docker-env)

for app in ${apps[@]}; do
    dockerfile_path="${app}/deployments/Dockerfile"
    if [[ -e $dockerfile_path ]]; then
        docker build . -t ${app}-app -f $dockerfile_path
    fi
done
