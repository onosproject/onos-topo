ARG ONOS_BUILD_VERSION=v0.5.0

FROM onosproject/golang-build:$ONOS_BUILD_VERSION
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-topo
RUN cd /go/src/github.com/onosproject/onos-topo && GOFLAGS=-mod=vendor make build
