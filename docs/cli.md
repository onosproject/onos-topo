<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
SPDX-License-Identifier: Apache-2.0
-->

# Command-Line Interface
The project provides a command-line facilities for remotely 
interacting with the topology subsystem.

The commands are available at run-time using the consolidated `onos` client hosted in the `onos-cli` repository.

The documentation about building and deploying the consolidated `onos` client or its Docker container
is available in the `onos-cli` GitHub repository.

## Usage
To see the detailed usage help for the `onos topo ...` family of commands, 
please see the [CLI documentation](https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_topo.md)

## Examples
Here are some concrete examples of usage:

List all entities.
```bash
$ onos topo get entities
Entity ID   Kind ID   Labels   Aspects
14550002    e2cell    <None>   onos.topo.E2Cell,onos.topo.Location,onos.topo.Coverage
14550001    e2cell    <None>   onos.topo.Coverage,onos.topo.E2Cell,onos.topo.Location
5154        e2node    <None>   onos.topo.E2Node
1454c003    e2cell    <None>   onos.topo.Location,onos.topo.Coverage,onos.topo.E2Cell
14550003    e2cell    <None>   onos.topo.E2Cell,onos.topo.Location,onos.topo.Coverage
5153        e2node    <None>   onos.topo.E2Node
1454c001    e2cell    <None>   onos.topo.Location,onos.topo.Coverage,onos.topo.E2Cell
1454c002    e2cell    <None>   onos.topo.E2Cell,onos.topo.Location,onos.topo.Coverage
```

List all entities of `e2node` kind.
```bash
$ onos topo get entities --kind e2node
Entity ID   Kind ID   Labels   Aspects
5153        e2node    <None>   onos.topo.E2Node
5154        e2node    <None>   onos.topo.E2Node
```
List all `e2cell` entities related to the specified `e2node` via `contains` relation.

```bash
$ onos topo get entities --related-to 5153 --related-via contains
1454c003   e2cell   <None>   onos.topo.E2Cell
1454c002   e2cell   <None>   onos.topo.E2Cell
1454c001   e2cell   <None>   onos.topo.E2Cell
```

Show verbose information on entity `1454c001`
```bash
$ onos topo get entity 1454c001 -v
1454c001   e2cell   <None>
           onos.topo.Location={"lat":52.486405,"lng":13.412234}
           onos.topo.Coverage={"arc_width":120,"azimuth":0,"height":43,"tilt":1}
           onos.topo.E2Cell={"cellObjectId":"13842601454c001","cellGlobalId":{"value":"1454c001"}}
```

Show all `neighbors` relations
```bash
$ onos topo get relations --kind neighbors
Relation ID         Kind ID     Source ID   Target ID   Labels   Aspects
1454c003-1454c002   neighbors   1454c003    1454c002    <None>   <None>
1454c001-1454c002   neighbors   1454c001    1454c002    <None>   <None>
1454c002-1454c003   neighbors   1454c002    1454c003    <None>   <None>
14550001-14550003   neighbors   14550001    14550003    <None>   <None>
14550002-14550001   neighbors   14550002    14550001    <None>   <None>
14550002-14550003   neighbors   14550002    14550003    <None>   <None>
1454c003-1454c001   neighbors   1454c003    1454c001    <None>   <None>
1454c001-1454c003   neighbors   1454c001    1454c003    <None>   <None>
1454c002-1454c001   neighbors   1454c002    1454c001    <None>   <None>
14550001-14550002   neighbors   14550001    14550002    <None>   <None>
14550003-14550002   neighbors   14550003    14550002    <None>   <None>
14550003-14550001   neighbors   14550003    14550001    <None>   <None>
```

Create a new entity with sparsely populated `Configurable` aspect
```bash
$ onos topo create entity "virtual" --aspect onos.topo.Configurable='{"type": "devicesim-1.0.x", "version": "1.0.0"}'
```
