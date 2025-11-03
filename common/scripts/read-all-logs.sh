#!/usr/bin/env bash

source $(dirname $0)/envs.sh

if [[ " $@ " =~ ( --minikube ) ]]; then
    minikube ssh -- lnav \
        $(printf "/tmp/hostpath-provisioner/azazon/%s-log-pvc/*.log " "${apps[@]}")
else 
    lnav $(printf "%s/logs/*.log " "${apps[@]}")
fi

