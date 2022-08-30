<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
SPDX-License-Identifier: Apache-2.0
-->

# onos-topo
[![Go Report Card](https://goreportcard.com/badge/github.com/onosproject/onos-topo)](https://goreportcard.com/report/github.com/onosproject/onos-topo)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/onosproject/onos-topo?status.svg)](https://godoc.org/github.com/onosproject/onos-topo)

## Overview

The µONOS Topology subsystem provides topology management for µONOS core services and applications.
The topology subsystem structures the topology information as a set of objects, which can be either 
`Entity`, `Relation` or a `Kind`.

* `Entity` objects are nodes in a graph and are generally intended to represent network devices, control entities, 
control domains, and so on.
* `Relation` objects are edges in a graph and are intended to represent various types of relations between two 
`Entity` objects, e.g. `contains`, `controls`, `implements`.
* `Kind` objects can be thought of as template or schema objects that represent an entity or a relation kind. 
Strictly speaking, `Entity` or `Relation` instances do not have to be associated with a `Kind`, 
but maintaining `Kind` associations can be used for schema validation and speeding up queries and is 
therefore highly encouraged.

## API
The `onos-topo` subsystem exposes the topology information via a [gRPC API] that supports the above abstractions.

### Unique ID
Each `Entity`, `Relation` and `Kind` objects has a unique identifier that can be used to directly look it up, 
update or delete it.

### Aspects
The `Entity` and `Relation` objects themselves carry only the essential information for identifying them, 
associating them with a particular kind and in case of `Relation`, for associating the two - source and target - 
`Entity` objects. Clearly, while this is necessary, it is not sufficient to allow the platform or applications to 
track other pertinent information about the entities and relations.

Since different use-cases or applications require tracking different information, and these may vary for different 
types of devices or network domains, the topology schema must be extensible to carry various aspects of information.
This is where the notion of `Aspect` comes in. An `Aspect` is a collection of structured information, modeled as a 
Protobuf message (although this is not strictly necessary), which is attached to any type of object; generally 
mostly an `Entity` or a `Relation`.

Each object carries a mapping of aspect type (`TypeURL`) and Protobuf `Any` message. For example, to track a geo-location 
of a network element, one can associate `onos.topo.Location` instance, populated with the appropriate longitude and
latitude with the `Entity` that represents that network element, with the `Location` being defined as follows:
```proto
// Geographical location; expected value type of "location" aspect
message Location {
    double lat = 1;
    double lng = 2;
}
```

Similarly, to track information about the cone of signal coverage for a radio-unit, one can attach `onos.topo.Coverage`
instance to an `Entity` representing the radio unit, with `Coverage` being defined as follows:
```proto
// Area of coverage; expected value type of "coverage" aspect
message Coverage {
    int32 height = 1;
    int32 arc_width = 2;
    int32 azimuth = 3;
    int32 tilt = 4;
}
```

The [current list of aspects](https://github.com/onosproject/onos-api/tree/master/proto/onos/topo) defined in `onos-api` includes the following:
* `onos.topo.Asset` - basic asset information for the device: model, HW, SW versions, serial number, etc.
* `onos.topo.Location` - geo location coordinates
* `onos.topo.Configurable` - info for devices that support configuration via gNMI
* `onos.topo.MastershipState` - for tracking mastership role
* `onos.topo.TLSOptoins` - TLS connection options
* `onos.topo.Protocols` - for tracking connectivity state of supported device control protocols
* `onos.topo.Coverage` - radio unit signal coverage cone information
* `onos.topo.E2Node` - information about an O-RAN E2 node
* `onos.topo.E2Cell` - information about an O-RAN E2 cell
* `onos.topo.AdHoc` - for tracking ad-hoc key/value string attributes (not labels)

The above are merely examples of aspects. Network control platforms and applications can supply their own depending on
the needs of a particular use-case.

### Labels
To assist in categorization of the topology objects, each object can carry a number of labels as meta-data. 
Each label carries a value. 

For example the `deployment` label can have `production` or `staging` or `testing` as values. Or similarly,
`tier` label can have `access`, `fabric` or `backhaul` as values to indicate the area of network where the entity
belongs.

### Filters
The topology API provides a `List` method to obtain a collection of objects. The caller can specify a number of
different filters to narrow the results. All topology objects will be returned if the request does not contain 
any filters.

* Type Filter - specifies which type(s) of objects - `Entity`, `Relation` or `Kind` - should be included.
* Kind Filter - specifies which kind(s) of objects should be included, e.g. `contains`, `controls`
* Labels Filter - specifies which label name/value(s) should be included, e.g. `tier=fabric`
* Relation Filter - specifies target entities related to a given source entity via a relation of a given kind
   
Support for other filters may be added in the future.

## Distribution
The topology subsystem is available as a [Docker] image and deployed with [Helm]. To build the Docker image,
run `make images`.

### Visualizer
To assist developers in visualizing the entities and relations tracked by `onos-topo`, a simple graphic visualization
tool is available. It can be run locally via:
```bash
# Requires 'kubectl port-forward deploy/onos-topo 5150' to forward topo gRPC API
> go run cmd/topo-visualizer/topo-visualizer.go --service-address localhost:5150
```
and then simply opening `http://localhost:5152` using your web browser of choice.

Alternatively, the visualizer can be run directly from the `onos-topo` docker container via:
```bash
# Requires 'kubectl port-forward deploy/onos-topo 5152' to forward visualizer HTTP/WS traffic
> k exec -it deploy/onos-topo -- /usr/local/bin/topo-visualizer
```

The visualizer uses `onos-topo` API to watch changes occurring on the topology and forwards
these changes via web-socket to the browser where it renders the various entities and relations using
a simple force layout graph. This allows the view to dynamically adjust to reflect the current topology state.

Clicking on nodes (entities) and links (relations) will show the full contents of the entity or relation 
as JSON structure. Nodes can be dragged around and the entire graph can be zoomed and panned within the viewport.

The visualizer is presently under active development.

## See Also
* [Deployment](docs/deployment.md)
* [CLI examples](docs/cli.md)
* [API examples (Golang)](docs/api-go.md)


[gRPC API]: https://github.com/onosproject/onos-api/blob/master/proto/onos/topo/topo.proto
[topology subcommands]: https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_topo.md
[Docker]: https://www.docker.com/
[Helm]: https://helm.sh
