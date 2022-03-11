<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
SPDX-License-Identifier: Apache-2.0
-->

# API Examples for Golang
The following are a few examples on how to use the [Golang API] generated from the `onos-topo` gRPC API.

All the examples assume that the following import is available:
```go
import "github.com/onosproject/onos-api/go/topo"
```

Similarly, all examples assume that the `TopoService` client has been obtained and that context `ctx` was 
either provided or created via code similar to the following:
```go
client := topo.CreateTopoClient(conn)
...
ctx := context.Background()
```

## Create an Entity
The following example shows how to create an `Entity` of `e2cell` kind, with a couple of separate aspects of 
information - `Location`, `Coverage` and `E2Cell`:
```go
cell := &topo.Object{
    ID:   topo.ID(cellID),
    Type: topo.Object_ENTITY,
    Obj: &topo.Object_Entity{
        Entity: &topo.Entity{
            KindID: topo.ID(topo.E2CELL),
        },
    }
}
cell.SetAspect(&topo.Location{Lat: 50.08834, Lng: 14.40372})
cell.SetAspect(&topo.Coverage{Height: 8.7, ArcWidth: 120.0, Azimuth: 315.0, Tilt: -5.0})
cell.SetAspect(&topo.E2Cell{CellObjectID: "4269A20", AntennaCount: 5, EARFCN: 69, PCI: 42, CellType: "FEMTO"})

resp, err := client.Create(ctx, &topo.CreateRequest{Object: cell})
```

## Create a Relation
Here we can see an example of creating `Relation` of `neighbors` kind, representing one cell being a neighbor 
of another. There are no aspects annotating this relation. Also, note that if the relation ID is unspecified 
during the creation, one will be automatically generated using the `topo.RelationID(...)` method, based on
the source, kind and the target IDs.
```go
relation := &topo.Object{
    Type: topo.Object_RELATION,
    Obj: &topo.Object_Relation{
        Entity: &topo.Relation{
            KindID: topo.ID(topo.NEIGHBORS),
            SrcEntityID: cellID,
            TgtEntityID: neighborCellID,
        },
    }
}

resp, err := client.Create(ctx, &topo.CreateRequest{Object: relation})
```

## Get an Object
In order to retrieve a specific object, one must simply provide its ID. To access a specific aspect, create
an instance of the aspect object and pass its reference to the `GetAspect` method of the topology object:
```go
cell, err := client.Get(ctx, &topo.GetRequest{ID: cellID})
if err == nil { ... }
coverage := &topo.Coverage{}
cell.GetAspect(coverage)
```

## Update an Entity
To update an aspect of an object, one must first obtain the object via `Get`, update its structure and
then call the `Update` method:
```go
cell, err := client.Get(&topo.GetRequest{ID: cellID})
if err == nil { ... }
cell.SetAspect(&topo.Location{Lat: 50.08834, Lng: 14.40372})

resp, err := client.Update(ctx, &topo.UpdateRequest{Object: cell})
```

## List Objects
The `List` method can be used to obtain various collections of objects using `Filters` specified as part of
the request. Presently, there are several types of filters.

* Object Type filter - specifies which types of objects are relevant
* Entity or Relation Kind filter - specifies which kinds of entities or relations are relevant
* Label filters - specifies criteria to filter on specific labels and their values
* Relation filter - allows finding objects related to another object via a particular kind of relation

The object kind and label filters can specify values via `equal`, `in` and `not` operators.

Furthermore, the results can be optionally sorted (ascending or descending) by the ID of the objects.
This can be useful for applications where determinism in the returned results is desirable.

Here is an example of requesting only `Entity` objects, of either `e2node` or `e2cell` kind, sorted in the
ascending order of their IDs:
```go
filters := &topo.Filters{
	KindFilter: &topoapi.Filter{
        Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: []string{topo.E2NODE, topo.E2CELL}}},
    },
}
resp, err := client.List(ctx, &topo.ListRequest{Filters: filters, SortOrder: topo.SortOrder_ASCENDING})
```

The following shows how to obtain an unordered list of `e2cell` entities, related to a specified node entity
using a `contains` relation:
```go
resp, err := client.List(ctx, &topo.ListRequest{Filters: &topo.Filters{
                RelationFilter: &topo.RelationFilter{
                    SrcId:        nodeID,
                    RelationKind: topo.CONTAINS,
                    TargetKind:   topo.E2CELL,
                },
            }})
```


## Watch Changes
The topology API allows clients to watch the changes in real-time via its `Watch` method which delivers its 
results as a continuous stream of events. These include not only the usual `create`, `update`, and `delete` events,
but also `replay` events to indicate the object as it existed prior to the `Watch` being called.

As with the `List` method, the results can be further narrowed by specifying `Filters` in the request.
Here is a simple example of the `Watch` usage, which does not specify any filters:

```go
stream, err := client.Watch(ctx, &topo.WatchRequest{})
if err == nil { ... }

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil { ... }
    processEvent(msg.Event.Type, msg.Event.Object)
}
```
The client can cancel the watch at anytime by invoking `ctx.Done()`.

## Delete an Object
Deleting an object requires to merely provide its ID:
```go
node, err := client.Delete(ctx, &topo.DeleteRequest{ID: nodeID})
if err == nil { ... }
```

[Golang API]: https://github.com/onosproject/onos-api/tree/master/go/onos/topo
