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
    - [AddEntityReq](#topo.AddEntityReq)
    - [AddEntityRespose](#topo.AddEntityRespose)
    - [GetAttrReq](#topo.GetAttrReq)
    - [GetAttrResp](#topo.GetAttrResp)
    - [GetAttrResp.AttrEntry](#topo.GetAttrResp.AttrEntry)
    - [RemoveEntityReq](#topo.RemoveEntityReq)
    - [RemoveEntityResp](#topo.RemoveEntityResp)
    - [UpdateEntityReq](#topo.UpdateEntityReq)
    - [UpdateEntityResp](#topo.UpdateEntityResp)
  
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



<a name="topo.AddEntityReq"></a>

### AddEntityReq
AddEntityReq adds a entity to the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  |  |
| id | [string](#string) |  |  |






<a name="topo.AddEntityRespose"></a>

### AddEntityRespose
AddEntityRespose is sent in response to an AddEntityReq






<a name="topo.GetAttrReq"></a>

### GetAttrReq
GetAttrReq ...


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  |  |
| id | [string](#string) |  |  |
| attrKind | [string](#string) |  |  |






<a name="topo.GetAttrResp"></a>

### GetAttrResp
GetAttrsResp carries a entity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  |  |
| id | [string](#string) |  |  |
| attr | [GetAttrResp.AttrEntry](#topo.GetAttrResp.AttrEntry) | repeated |  |






<a name="topo.GetAttrResp.AttrEntry"></a>

### GetAttrResp.AttrEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.RemoveEntityReq"></a>

### RemoveEntityReq
RemoveEntityReq removes a entity by ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  |  |
| id | [string](#string) |  |  |






<a name="topo.RemoveEntityResp"></a>

### RemoveEntityResp
RemoveEntityResp is sent in response to a RemoveEntityReq






<a name="topo.UpdateEntityReq"></a>

### UpdateEntityReq
UpdateEntityReq updates a entity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  |  |
| id | [string](#string) |  |  |






<a name="topo.UpdateEntityResp"></a>

### UpdateEntityResp
UpdateEntityResp is sent in response to an UpdateEntityReq





 

 

 


<a name="topo.entityService"></a>

### entityService
EntityService provides an API for managing entities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Add | [AddEntityReq](#topo.AddEntityReq) | [AddEntityRespose](#topo.AddEntityRespose) | Add adds a entity to the topology |
| Update | [UpdateEntityReq](#topo.UpdateEntityReq) | [UpdateEntityResp](#topo.UpdateEntityResp) | Update updates a entity |
| GetAttr | [GetAttrReq](#topo.GetAttrReq) | [GetAttrResp](#topo.GetAttrResp) | Get gets a entity by ID |
| RemoveEntity | [RemoveEntityReq](#topo.RemoveEntityReq) | [RemoveEntityResp](#topo.RemoveEntityResp) | Remove removes a entity from the topology |

 



<a name="api/topo/interface.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/topo/interface.proto



<a name="topo.Interface"></a>

### Interface
Interface contains information about a interface






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

