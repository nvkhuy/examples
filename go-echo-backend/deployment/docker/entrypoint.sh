#!/bin/bash

jq -S -nr env | jq '{SERVICE_NAME: .SERVICE_NAME, SERVICE_NAME: .SERVICE_NAME, GIT_VERSION: .GIT_VERSION, BUILD_TIMESTAMP: .BUILD_TIMESTAMP}'

cp /app/efs/$ENV.json /app/$ENV.json

exec "$@"
