FROM debian:bullseye
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV GO_VERSION=1.22.1
COPY . /var/build
RUN echo '[INFO] --- update repo index' && \
    apt-get update && \
    echo '[INFO] --- install essential packages' && \
    apt-get install -y wget && \
    echo '[INFO] --- install go' && \
    cd /tmp && \
    wget  --no-verbose "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local/ -xzf "go${GO_VERSION}.linux-amd64.tar.gz" || true 

RUN echo '[INFO] --- build go project' && \
    cd /var/build && \
    pwd && \
    find && \
    ./scripts/build.sh "${BUILD_TAG}"
