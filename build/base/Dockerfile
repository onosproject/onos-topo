ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:$ONOS_BUILD_VERSION
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-topo
RUN cd /go/src/github.com/onosproject/onos-topo && GOFLAGS=-mod=vendor make build
