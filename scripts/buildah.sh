#!/bin/bash

# argv
APP_VERSION=$1
if [ -z "${APP_VERSION}" ]; then
    echo "Usage: $0 APP_VERSION"
    exit 1
fi

mkdir -p ./build
buildah rmi builder &>/dev/null
buildah build \
    -t builder \
    -v "${PWD}/build":"/var/work/build" \
    -f Containerfile \
    --env APP_VERSION="${APP_VERSION}" \
    .
