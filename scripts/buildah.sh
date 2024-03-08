mkdir -p ./build
buildah rmi builder &>/dev/null
buildah build \
    -t builder \
    -v "${PWD}":"/var/build/build" \
    -f Containerfile \
    .
