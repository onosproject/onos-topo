# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/topo/topo.proto](#api/topo/topo.proto)
    - [CreateRequest](#topo.CreateRequest)
    - [CreateResponse](#topo.CreateResponse)
    - [DeleteRequest](#topo.DeleteRequest)
    - [DeleteResponse](#topo.DeleteResponse)
    - [Entity](#topo.Entity)
    - [Event](#topo.Event)
    - [GetRequest](#topo.GetRequest)
    - [GetResponse](#topo.GetResponse)
    - [Kind](#topo.Kind)
    - [Kind.AttributesEntry](#topo.Kind.AttributesEntry)
    - [ListRequest](#topo.ListRequest)
    - [ListResponse](#topo.ListResponse)
    - [Object](#topo.Object)
    - [Object.AttributesEntry](#topo.Object.AttributesEntry)
    - [Relation](#topo.Relation)
    - [UpdateRequest](#topo.UpdateRequest)
    - [UpdateResponse](#topo.UpdateResponse)
    - [WatchRequest](#topo.WatchRequest)
    - [WatchResponse](#topo.WatchResponse)
  
    - [EventType](#topo.EventType)
    - [Object.Type](#topo.Object.Type)
  
    - [Topo](#topo.Topo)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/topo/topo.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/topo.proto



<a name="topo.CreateRequest"></a>

### CreateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.CreateResponse"></a>

### CreateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.DeleteRequest"></a>

### DeleteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="topo.DeleteResponse"></a>

### DeleteResponse







<a name="topo.Entity"></a>

### Entity
Entity represents any &#34;thing&#34; that is represented in the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind_id | [string](#string) |  | user-defined entity kind |
| protocols | [device.ProtocolState](#topo.device.ProtocolState) | repeated |  |






<a name="topo.Event"></a>

### Event
Event is a topo operation event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [EventType](#topo.EventType) |  |  |
| object | [Object](#topo.Object) |  |  |






<a name="topo.GetRequest"></a>

### GetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="topo.GetResponse"></a>

### GetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object](#topo.Object) |  |  |






<a name="topo.Kind"></a>

### Kind



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| attributes | [Kind.AttributesEntry](#topo.Kind.AttributesEntry) | repeated | Map of attributes and their default values for this Kind |






<a name="topo.Kind.AttributesEntry"></a>

### Kind.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.ListRequest"></a>

### ListRequest







<a name="topo.ListResponse"></a>

### ListResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.Object"></a>

### Object



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| revision | [uint64](#uint64) |  |  |
| type | [Object.Type](#topo.Object.Type) |  |  |
| entity | [Entity](#topo.Entity) |  |  |
| relation | [Relation](#topo.Relation) |  |  |
| kind | [Kind](#topo.Kind) |  |  |
| attributes | [Object.AttributesEntry](#topo.Object.AttributesEntry) | repeated |  |






<a name="topo.Object.AttributesEntry"></a>

### Object.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.Relation"></a>

### Relation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind_id | [string](#string) |  | user defined relation kind |
| src_entity_id | [string](#string) |  |  |
| tgt_entity_id | [string](#string) |  |  |






<a name="topo.UpdateRequest"></a>

### UpdateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.UpdateResponse"></a>

### UpdateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.WatchRequest"></a>

### WatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| noreplay | [bool](#bool) |  |  |






<a name="topo.WatchResponse"></a>

### WatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [Event](#topo.Event) |  |  |





 


<a name="topo.EventType"></a>

### EventType
EventType is a topo operation event type

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| ADDED | 1 |  |
| UPDATED | 2 |  |
| REMOVED | 3 |  |



<a name="topo.Object.Type"></a>

### Object.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 |  |
| ENTITY | 1 |  |
| RELATION | 2 |  |
| KIND | 3 |  |


 

 


<a name="topo.Topo"></a>

### Topo
EntityService provides an API for managing entities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateRequest](#topo.CreateRequest) | [CreateResponse](#topo.CreateResponse) | Create a new topology object |
| Get | [GetRequest](#topo.GetRequest) | [GetResponse](#topo.GetResponse) | Get an object from topology |
| Update | [UpdateRequest](#topo.UpdateRequest) | [UpdateResponse](#topo.UpdateResponse) | Update an existing topology object |
| Delete | [DeleteRequest](#topo.DeleteRequest) | [DeleteResponse](#topo.DeleteResponse) | Delete an object from topology |
| List | [ListRequest](#topo.ListRequest) | [ListResponse](#topo.ListResponse) | List gets a stream of requested objects |
| Watch | [WatchRequest](#topo.WatchRequest) | [WatchResponse](#topo.WatchResponse) stream | Watch returns a stream of topo change notifications |

 



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

