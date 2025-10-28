#!/usr/bin/env bash

if [[ " $@ " =~ ( --minikube ) ]]; then
    minikube ssh -- lnav \
        /tmp/hostpath-provisioner/azazon/{auth,notification,order,payment,stock}-log-pvc/*.log
else 
    lnav {auth,notification,order,payment,stock}/logs/*.log
fi

