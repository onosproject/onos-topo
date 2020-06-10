# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/topo/topo.proto](#api/topo/topo.proto)
    - [Entity](#topo.Entity)
    - [Entity.E2Interface](#topo.Entity.E2Interface)
    - [Entity.E2Node](#topo.Entity.E2Node)
    - [Entity.Ric](#topo.Entity.Ric)
    - [Entity.XnInterface](#topo.Entity.XnInterface)
    - [ReadRequest](#topo.ReadRequest)
    - [ReadResponse](#topo.ReadResponse)
    - [Relationship](#topo.Relationship)
    - [Relationship.Aggregates](#topo.Relationship.Aggregates)
    - [Relationship.Contains](#topo.Relationship.Contains)
    - [StreamMessageRequest](#topo.StreamMessageRequest)
    - [StreamMessageResponse](#topo.StreamMessageResponse)
    - [Update](#topo.Update)
    - [WriteRequest](#topo.WriteRequest)
    - [WriteResponse](#topo.WriteResponse)
  
    - [Entity.Kind](#topo.Entity.Kind)
    - [Relationship.Kind](#topo.Relationship.Kind)
    - [Update.Type](#topo.Update.Type)
  
    - [topo](#topo.topo)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/topo/topo.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/topo.proto



<a name="topo.Entity"></a>

### Entity
Entity represent &#34;things&#34;


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [Entity.Kind](#topo.Entity.Kind) |  |  |
| ric | [Entity.Ric](#topo.Entity.Ric) |  |  |
| e2_node | [Entity.E2Node](#topo.Entity.E2Node) |  |  |
| e2_interface | [Entity.E2Interface](#topo.Entity.E2Interface) |  |  |
| xn_interface | [Entity.XnInterface](#topo.Entity.XnInterface) |  |  |
| relationships | [Relationship](#topo.Relationship) | repeated |  |






<a name="topo.Entity.E2Interface"></a>

### Entity.E2Interface







<a name="topo.Entity.E2Node"></a>

### Entity.E2Node







<a name="topo.Entity.Ric"></a>

### Entity.Ric







<a name="topo.Entity.XnInterface"></a>

### Entity.XnInterface







<a name="topo.ReadRequest"></a>

### ReadRequest







<a name="topo.ReadResponse"></a>

### ReadResponse







<a name="topo.Relationship"></a>

### Relationship



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [Relationship.Kind](#topo.Relationship.Kind) |  |  |
| contains | [Relationship.Contains](#topo.Relationship.Contains) |  |  |
| aggregates | [Relationship.Aggregates](#topo.Relationship.Aggregates) |  |  |






<a name="topo.Relationship.Aggregates"></a>

### Relationship.Aggregates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AggregatorId | [string](#string) |  |  |
| AggregateeId | [string](#string) |  |  |






<a name="topo.Relationship.Contains"></a>

### Relationship.Contains



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ContainerId | [string](#string) |  |  |
| ContaineeId | [string](#string) |  |  |






<a name="topo.StreamMessageRequest"></a>

### StreamMessageRequest







<a name="topo.StreamMessageResponse"></a>

### StreamMessageResponse







<a name="topo.Update"></a>

### Update



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Update.Type](#topo.Update.Type) |  |  |
| entity | [Entity](#topo.Entity) |  |  |






<a name="topo.WriteRequest"></a>

### WriteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| updates | [Update](#topo.Update) | repeated | The write batch, comprising a list of Update operations |






<a name="topo.WriteResponse"></a>

### WriteResponse






 


<a name="topo.Entity.Kind"></a>

### Entity.Kind


| Name | Number | Description |
| ---- | ------ | ----------- |
| RIC | 0 |  |
| E2NODE | 1 |  |
| E2INTERFACE | 2 |  |
| XNINTERFACE | 3 |  |



<a name="topo.Relationship.Kind"></a>

### Relationship.Kind


| Name | Number | Description |
| ---- | ------ | ----------- |
| CONTAINS | 0 |  |
| AGGREGATES | 1 |  |



<a name="topo.Update.Type"></a>

### Update.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 |  |
| INSERT | 1 |  |
| MODIFY | 2 |  |
| DELETE | 3 |  |


 

 


<a name="topo.topo"></a>

### topo
EntityService provides an API for managing entities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Write | [WriteRequest](#topo.WriteRequest) | [WriteResponse](#topo.WriteResponse) | Update one or more entities to the topology |
| Read | [ReadRequest](#topo.ReadRequest) | [ReadResponse](#topo.ReadResponse) | Read one or more entities from topology |
| StreamChannel | [StreamMessageRequest](#topo.StreamMessageRequest) stream | [StreamMessageResponse](#topo.StreamMessageResponse) stream | Represents the bidirectional stream between onos-topo and a client for the purpose of - streaming notifications |

 



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

