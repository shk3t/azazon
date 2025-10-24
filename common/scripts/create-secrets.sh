#!/usr/bin/env bash

kubectl create namespace security

kubectl create secret tls temp-tls-secret \
    --cert=cert/tls.crt \
    --key=cert/tls.key \
    --namespace=security
