# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/topo/topo.proto](#api/topo/topo.proto)
    - [Attributes](#topo.Attributes)
    - [Attributes.AttrsEntry](#topo.Attributes.AttrsEntry)
    - [Entity](#topo.Entity)
    - [Object](#topo.Object)
    - [ReadRequest](#topo.ReadRequest)
    - [ReadResponse](#topo.ReadResponse)
    - [Reference](#topo.Reference)
    - [Relationship](#topo.Relationship)
    - [SubscribeRequest](#topo.SubscribeRequest)
    - [SubscribeResponse](#topo.SubscribeResponse)
    - [Update](#topo.Update)
    - [WriteRequest](#topo.WriteRequest)
    - [WriteResponse](#topo.WriteResponse)
  
    - [Object.Type](#topo.Object.Type)
    - [Relationship.Directionality](#topo.Relationship.Directionality)
    - [Relationship.Multiplicity](#topo.Relationship.Multiplicity)
    - [Relationship.Type](#topo.Relationship.Type)
    - [Update.Type](#topo.Update.Type)
  
    - [Topo](#topo.Topo)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/topo/topo.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/topo.proto



<a name="topo.Attributes"></a>

### Attributes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attrs | [Attributes.AttrsEntry](#topo.Attributes.AttrsEntry) | repeated | TODO - Instead of a plain string, consider using a &#34;typed&#34; value in attrs map. - See onos-config for example. |






<a name="topo.Attributes.AttrsEntry"></a>

### Attributes.AttrsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.Entity"></a>

### Entity
Entity represents any &#34;thing&#34; that is represented in the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  |  |






<a name="topo.Object"></a>

### Object



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ref | [Reference](#topo.Reference) |  |  |
| type | [Object.Type](#topo.Object.Type) |  |  |
| entity | [Entity](#topo.Entity) |  |  |
| relationship | [Relationship](#topo.Relationship) |  |  |
| attrs | [Attributes](#topo.Attributes) |  |  |






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






<a name="topo.Relationship"></a>

### Relationship



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| directionality | [Relationship.Directionality](#topo.Relationship.Directionality) |  |  |
| multiplicity | [Relationship.Multiplicity](#topo.Relationship.Multiplicity) |  |  |
| type | [Relationship.Type](#topo.Relationship.Type) |  |  |
| source_refs | [Reference](#topo.Reference) | repeated | The two sets of objects that the relationship binds |
| target_refs | [Reference](#topo.Reference) | repeated |  |






<a name="topo.SubscribeRequest"></a>

### SubscribeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ref | [Reference](#topo.Reference) |  |  |
| withoutReplay | [bool](#bool) |  |  |






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
| RELATIONSHIP | 2 |  |



<a name="topo.Relationship.Directionality"></a>

### Relationship.Directionality


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED_DIRECTIONALITY | 0 |  |
| DIRECTED | 1 |  |
| BIDIRECTIONAL | 2 |  |



<a name="topo.Relationship.Multiplicity"></a>

### Relationship.Multiplicity


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED_MULTIPLICITY | 0 |  |
| ONE_TO_ONE | 1 |  |
| ONE_TO_MANY | 2 |  |
| MANY_TO_ONE | 3 |  |
| MANY_TO_MANY | 4 |  |



<a name="topo.Relationship.Type"></a>

### Relationship.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 |  |
| CONTAINS | 1 |  |
| CONTROLS | 2 |  |
| AGGREGATES | 3 |  |
| ORIGINATES | 4 |  |
| TERMINATES | 5 |  |
| TRAVERSES | 6 |  |
| REALIZED_BY | 7 |  |



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

