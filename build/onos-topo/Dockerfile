ARG ONOS_TOPO_BASE_VERSION=latest

FROM onosproject/onos-topo-base:$ONOS_TOPO_BASE_VERSION as base

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-topo/build/_output/onos-topo /usr/local/bin/onos-topo

ENTRYPOINT ["onos-topo"]
