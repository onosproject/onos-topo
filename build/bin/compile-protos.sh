#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"

protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=topo/admin,plugins=grpc:. pkg/northbound/admin/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=topo/device,plugins=grpc:. pkg/northbound/device/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=topo/diags,plugins=grpc:. pkg/northbound/diags/*.proto
