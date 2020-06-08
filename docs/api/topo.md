# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/topo/e2interface.proto](#api/topo/e2interface.proto)
    - [E2Interface](#topo.E2Interface)
  
- [api/topo/e2node.proto](#api/topo/e2node.proto)
    - [E2Node](#topo.E2Node)
    - [E2Node.Interfaces](#topo.E2Node.Interfaces)
  
- [api/topo/endpoint.proto](#api/topo/endpoint.proto)
    - [EndPoint](#topo.EndPoint)
  
- [api/topo/entity.proto](#api/topo/entity.proto)
    - [AddRequest](#topo.AddRequest)
    - [AddResponse](#topo.AddResponse)
    - [Entity](#topo.Entity)
    - [Entity.ContainsEntry](#topo.Entity.ContainsEntry)
    - [GetRequest](#topo.GetRequest)
    - [GetResponse](#topo.GetResponse)
    - [ListRequest](#topo.ListRequest)
    - [ListResponse](#topo.ListResponse)
    - [RemoveRequest](#topo.RemoveRequest)
    - [RemoveResponse](#topo.RemoveResponse)
    - [UpdateRequest](#topo.UpdateRequest)
    - [UpdateResponse](#topo.UpdateResponse)
  
    - [ListResponse.Type](#topo.ListResponse.Type)
  
    - [entityService](#topo.entityService)
  
- [api/topo/interface.proto](#api/topo/interface.proto)
    - [Interface](#topo.Interface)
    - [Interfaces](#topo.Interfaces)
  
- [api/topo/link.proto](#api/topo/link.proto)
    - [Link](#topo.Link)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/topo/e2interface.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/e2interface.proto



<a name="topo.E2Interface"></a>

### E2Interface
E2Interface contains information about a E2 Interface





 

 

 

 



<a name="api/topo/e2node.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/e2node.proto



<a name="topo.E2Node"></a>

### E2Node
E2Node contains information about a E2 Node device






<a name="topo.E2Node.Interfaces"></a>

### E2Node.Interfaces



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| e2 | [E2Interface](#topo.E2Interface) | repeated |  |





 

 

 

 



<a name="api/topo/endpoint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/endpoint.proto



<a name="topo.EndPoint"></a>

### EndPoint
EndPoint represents the endpoint of a link


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| link | [Link](#topo.Link) |  |  |





 

 

 

 



<a name="api/topo/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/entity.proto



<a name="topo.AddRequest"></a>

### AddRequest
AddRequest adds a entity to the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the entity to add |






<a name="topo.AddResponse"></a>

### AddResponse
AddResponse is sent in response to an AddEntityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the entity with a revision number |






<a name="topo.Entity"></a>

### Entity
Entity represents &#34;things&#34; in a network topology like devices or links


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is a globally unique entity identifier |
| type | [string](#string) |  | type is the type of the entity |
| contains | [Entity.ContainsEntry](#topo.Entity.ContainsEntry) | repeated |  |






<a name="topo.Entity.ContainsEntry"></a>

### Entity.ContainsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.GetRequest"></a>

### GetRequest
GetRequest gets a entity by ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is the unique entity ID with which to lookup the entity |






<a name="topo.GetResponse"></a>

### GetResponse
GetResponse carries a entity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the entity object |






<a name="topo.ListRequest"></a>

### ListRequest
ListRequest requests a stream of entities and changes
By default, the request requests a stream of all entities that are present in the topology when
the request is received by the service. However, if `subscribe` is `true`, the stream will remain
open after all entities have been sent and events that occur following the last entity will be
streamed to the client until the stream is closed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscribe | [bool](#bool) |  | subscribe indicates whether to subscribe to events (e.g. ADD, UPDATE, and REMOVE) that occur after all entities have been streamed to the client |






<a name="topo.ListResponse"></a>

### ListResponse
ListResponse carries a single entity event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ListResponse.Type](#topo.ListResponse.Type) |  | type is the type of the event |
| entity | [Entity](#topo.Entity) |  | entity is the entity on which the event occurred |






<a name="topo.RemoveRequest"></a>

### RemoveRequest
RemoveRequest removes a entity by ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the entity to remove |






<a name="topo.RemoveResponse"></a>

### RemoveResponse
RemoveResponse is sent in response to a RemoveEntityRequest






<a name="topo.UpdateRequest"></a>

### UpdateRequest
UpdateRequest updates a entity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the updated entity |






<a name="topo.UpdateResponse"></a>

### UpdateResponse
UpdateResponse is sent in response to an UpdateEntityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entity | [Entity](#topo.Entity) |  | entity is the entity with updated revision |





 


<a name="topo.ListResponse.Type"></a>

### ListResponse.Type
Entity event type

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 | NONE indicates this response does not represent a state change |
| ADDED | 1 | ADDED is an event which occurs when a entity is added to the topology |
| UPDATED | 2 | UPDATED is an event which occurs when a entity is updated |
| REMOVED | 3 | REMOVED is an event which occurs when a entity is removed from the topology |


 

 


<a name="topo.entityService"></a>

### entityService
EntityService provides an API for managing entities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Add | [AddRequest](#topo.AddRequest) | [AddResponse](#topo.AddResponse) | Add adds a entity to the topology |
| Update | [UpdateRequest](#topo.UpdateRequest) | [UpdateResponse](#topo.UpdateResponse) | Update updates a entity |
| Get | [GetRequest](#topo.GetRequest) | [GetResponse](#topo.GetResponse) | Get gets a entity by ID |
| List | [ListRequest](#topo.ListRequest) | [ListResponse](#topo.ListResponse) stream | List gets a stream of entity add/update/remove events |
| Remove | [RemoveRequest](#topo.RemoveRequest) | [RemoveResponse](#topo.RemoveResponse) | Remove removes a entity from the topology |

 



<a name="api/topo/interface.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/interface.proto



<a name="topo.Interface"></a>

### Interface
Interface contains information about a interface


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| iface | [google.protobuf.Any](#google.protobuf.Any) |  | device is the actual device type (e.g. E2 Node, RIC) |






<a name="topo.Interfaces"></a>

### Interfaces



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ifaces | [Interface](#topo.Interface) | repeated |  |





 

 

 

 



<a name="api/topo/link.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/link.proto



<a name="topo.Link"></a>

### Link
Interface contains information about a interface





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

