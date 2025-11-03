#!/usr/bin/env bash

kubectl get pv | grep Released | awk '{print $1}' | xargs kubectl delete pv
