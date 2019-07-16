ARG ONOS_TOPO_BASE_VERSION=latest

FROM onosproject/onos-topo-base:$ONOS_TOPO_BASE_VERSION as base

FROM golang:1.12.6-alpine3.9 as debugBuilder

RUN apk upgrade --update --no-cache && apk add git && \
    go get -u github.com/go-delve/delve/cmd/dlv && \
    go get -u github.com/atomix/atomix-cli/cmd/atomix

FROM alpine:3.9

RUN apk upgrade --update --no-cache && apk add bash bash-completion libc6-compat

COPY --from=base /go/src/github.com/onosproject/onos-topo/build/_output/onos /usr/local/bin/onos
COPY --from=base /go/src/github.com/onosproject/onos-topo/build/_output/onos-topo-debug /usr/local/bin/onos-topo
COPY --from=debugBuilder /go/bin/dlv /usr/local/bin/dlv
COPY --from=debugBuilder /go/bin/atomix /usr/local/bin/atomix

RUN echo "#!/bin/sh" >> /usr/local/bin/onos-topo-debug && \
    echo "atomix topoler set \$ATOMIX_TOPOLER" >> /usr/local/bin/onos-topo-debug && \
    echo "atomix config set namespace \$ATOMIX_NAMESPACE" >> /usr/local/bin/onos-topo-debug && \
    echo "atomix config set app \$ATOMIX_APP" >> /usr/local/bin/onos-topo-debug && \
    echo "dlv --listen=:40000 --headless=true --accept-multiclient=true --api-version=2 exec /usr/local/bin/onos-topo -- \"\$@\"" >> /usr/local/bin/onos-topo-debug && \
    chmod +x /usr/local/bin/onos-topo-debug

RUN addgroup -S onos-topo && adduser -S -G onos-topo onos-topo
USER onos-topo
WORKDIR /home/onos-topo

ENV ATOMIX_TOPOLER=atomix-topoler.kube-system.svc.cluster.local:5679
ENV ATOMIX_NAMESPACE=default

RUN onos init && \
    cp /etc/profile /home/onos-topo/.bashrc && \
    onos completion bash > /home/onos-topo/.onos/bash_completion.sh && \
    echo "source /home/onos-topo/.onos/bash_completion.sh" >> /home/onos-topo/.bashrc && \
    onos config set address 127.0.0.1:5150 && \
    atomix init && \
    atomix completion bash > /home/onos-topo/.atomix/bash_completion.sh && \
    echo "source /home/onos-topo/.atomix/bash_completion.sh" >> /home/onos-topo/.bashrc

ENTRYPOINT ["onos-topo-debug"]
