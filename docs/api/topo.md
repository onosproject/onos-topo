# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/topo/topo.proto](#api/topo/topo.proto)
    - [Entity](#topo.Entity)
    - [Kind](#topo.Kind)
    - [Object](#topo.Object)
    - [Object.AttributesEntry](#topo.Object.AttributesEntry)
    - [ReadRequest](#topo.ReadRequest)
    - [ReadResponse](#topo.ReadResponse)
    - [Reference](#topo.Reference)
    - [Relation](#topo.Relation)
    - [SubscribeRequest](#topo.SubscribeRequest)
    - [SubscribeResponse](#topo.SubscribeResponse)
    - [Update](#topo.Update)
    - [WriteRequest](#topo.WriteRequest)
    - [WriteResponse](#topo.WriteResponse)
  
    - [Object.Type](#topo.Object.Type)
    - [Update.Type](#topo.Update.Type)
  
    - [Topo](#topo.Topo)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/topo/topo.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/topo.proto



<a name="topo.Entity"></a>

### Entity
Entity represents any &#34;thing&#34; that is represented in the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [Reference](#topo.Reference) |  | user-defined entity kind |






<a name="topo.Kind"></a>

### Kind



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| required_attrs | [string](#string) | repeated |  |






<a name="topo.Object"></a>

### Object



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ref | [Reference](#topo.Reference) |  |  |
| type | [Object.Type](#topo.Object.Type) |  |  |
| entity | [Entity](#topo.Entity) |  |  |
| relation | [Relation](#topo.Relation) |  |  |
| kind | [Kind](#topo.Kind) |  |  |
| attributes | [Object.AttributesEntry](#topo.Object.AttributesEntry) | repeated | TODO - Instead of a plain string, consider using a &#34;typed&#34; value in attributes map. - See onos-config for example. |






<a name="topo.Object.AttributesEntry"></a>

### Object.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.ReadRequest"></a>

### ReadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refs | [Reference](#topo.Reference) | repeated |  |






<a name="topo.ReadResponse"></a>

### ReadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [Object](#topo.Object) | repeated |  |






<a name="topo.Reference"></a>

### Reference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="topo.Relation"></a>

### Relation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [Reference](#topo.Reference) |  | user defined relation kind |
| source_ref | [Reference](#topo.Reference) |  |  |
| target_ref | [Reference](#topo.Reference) |  |  |






<a name="topo.SubscribeRequest"></a>

### SubscribeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ref | [Reference](#topo.Reference) |  |  |
| noreplay | [bool](#bool) |  |  |
| snapshot | [bool](#bool) |  |  |






<a name="topo.SubscribeResponse"></a>

### SubscribeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updates | [Update](#topo.Update) | repeated |  |






<a name="topo.Update"></a>

### Update



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Update.Type](#topo.Update.Type) |  |  |
| object | [Object](#topo.Object) |  |  |






<a name="topo.WriteRequest"></a>

### WriteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updates | [Update](#topo.Update) | repeated | The write batch, comprising a list of Update operations |






<a name="topo.WriteResponse"></a>

### WriteResponse






 


<a name="topo.Object.Type"></a>

### Object.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 |  |
| ENTITY | 1 |  |
| RELATION | 2 |  |
| KIND | 3 |  |



<a name="topo.Update.Type"></a>

### Update.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 |  |
| INSERT | 1 |  |
| MODIFY | 2 |  |
| DELETE | 3 |  |


 

 


<a name="topo.Topo"></a>

### Topo
EntityService provides an API for managing entities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Write | [WriteRequest](#topo.WriteRequest) | [WriteResponse](#topo.WriteResponse) | Update one or more entities to the topology |
| Read | [ReadRequest](#topo.ReadRequest) | [ReadResponse](#topo.ReadResponse) | Read one or more entities from topology |
| Subscribe | [SubscribeRequest](#topo.SubscribeRequest) | [SubscribeResponse](#topo.SubscribeResponse) stream | Subscribe returns a stream of topo change notifications |

 



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

