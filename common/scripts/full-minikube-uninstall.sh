#!/usr/bin/env bash

scriptpath=$(dirname $0)
source $scriptpath/envs.sh

helm uninstall -n $namespace $baseapp --ignore-not-found
for app in ${apps[@]}; do
    helm uninstall -n $namespace $app --ignore-not-found
    helm uninstall -n $namespace ${app}-db --ignore-not-found
done

eval $(minikube -p minikube docker-env)
docker image prune -f
for app in ${apps[@]}; do
    docker image remove -f ${app}-app
done

if [[ " $@ " =~ ( --drop-volumes ) ]]; then
    kubectl delete persistentvolumeclaims -n $namespace "data-0-${kafka_cluster_name}-${kafka_nodepool_name}-0" 2> /dev/null
    for app in ${apps[@]}; do
        kubectl delete persistentvolumeclaims -n $namespace -l app=${app}-db 2> /dev/null
    done
fi

bash $scriptpath/kube-delete-released-volumes.sh 2> /dev/null
