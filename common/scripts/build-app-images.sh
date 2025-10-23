#!/usr/bin/env bash

apps=(auth notification order payment stock)

eval $(minikube -p minikube docker-env)

for app in ${apps[@]}; do
    dockerfile_path="${app}/deployments/Dockerfile"
    if [[ -e $dockerfile_path ]]; then
        docker build . -t ${app}-app -f $dockerfile_path
    fi
done
