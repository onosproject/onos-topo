ARG ONOS_TOPO_BASE_VERSION=latest

FROM onosproject/onos-topo-base:$ONOS_TOPO_BASE_VERSION as base

FROM golang:1.12.6-alpine3.9 as debugBuilder

RUN apk upgrade --update --no-cache && apk add git && \
    go get -u github.com/go-delve/delve/cmd/dlv && \
    go get -u github.com/onosproject/onos-cli/cmd/onos

FROM alpine:3.9

RUN apk upgrade --update --no-cache && apk add bash bash-completion libc6-compat

COPY --from=base /go/src/github.com/onosproject/onos-topo/build/_output/onos-topo-debug /usr/local/bin/onos-topo
COPY --from=debugBuilder /go/bin/onos /usr/local/bin/onos
COPY --from=debugBuilder /go/bin/dlv /usr/local/bin/dlv

RUN echo "#!/bin/sh" >> /usr/local/bin/onos-topo-debug && \
    echo "dlv --listen=:40000 --headless=true --accept-multiclient=true --api-version=2 exec /usr/local/bin/onos-topo -- \"\$@\"" >> /usr/local/bin/onos-topo-debug && \
    chmod +x /usr/local/bin/onos-topo-debug

RUN addgroup -S onos-topo && adduser -S -G onos-topo onos-topo
USER onos-topo
WORKDIR /home/onos-topo

ENTRYPOINT ["onos-topo-debug"]
