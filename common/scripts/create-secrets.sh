#!/usr/bin/env bash

# mkdir cert && cd cert
# openssl genrsa -out tls.key 2048
# openssl req -new -key tls.key -out tls.csr -subj "/CN=azazon.com"
# openssl x509 -req -days 365 -in tls.csr -signkey tls.key -out tls.crt

kubectl create namespace security

kubectl create secret tls temp-tls-secret \
    --cert=cert/tls.crt \
    --key=cert/tls.key \
    --namespace=security
