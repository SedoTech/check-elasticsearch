#!/bin/bash

# argv
APP_VERSION=$1
if [ -z "${APP_VERSION}" ]; then
    APP_VERSION=development
fi

mkdir -p ./build/
go build -ldflags "-X main.version=${APP_VERSION}" -o ./build/check-elasticsearch && {
    chmod +x ./build/check-elasticsearch
}
