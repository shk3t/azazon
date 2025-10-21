#!/usr/bin/env bash

set -a
source .env
set +a

APP_DIR="notification/deployments/notification"
envsubst < $APP_DIR/values.yaml.tpl > $APP_DIR/values.yaml
