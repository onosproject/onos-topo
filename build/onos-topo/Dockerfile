# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

FROM onosproject/golang-build:v1.3.0 as build

ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-topo
COPY build/ /go/src/github.com/onosproject/onos-topo/build
RUN cd /go/src/github.com/onosproject/onos-topo && GOFLAGS=-mod=vendor make build

FROM alpine:3.16
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-topo/build/_output/onos-topo /usr/local/bin/onos-topo
COPY --from=build /go/src/github.com/onosproject/onos-topo/build/_output/topo-visualizer /usr/local/bin/topo-visualizer
COPY --from=build /go/src/github.com/onosproject/onos-topo/pkg/tools/topo-visualizer/index.html /var/topo-visualizer/index.html

ENTRYPOINT ["onos-topo"]
