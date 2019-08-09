#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/google/protobuf/src:${GOPATH}/src"

protoc -I=$proto_imports --go_out=import_path=topo/admin,plugins=grpc:. pkg/northbound/admin/*.proto
protoc -I=$proto_imports --go_out=import_path=topo/device,plugins=grpc:. pkg/northbound/device/*.proto
protoc -I=$proto_imports --go_out=import_path=topo/diags,plugins=grpc:. pkg/northbound/diags/*.proto
