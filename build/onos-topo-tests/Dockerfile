FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/onos-topo-tests/_output/bin/onos-topo-tests /usr/local/bin/onos-topo-tests

ENTRYPOINT ["onos-topo-tests"]
