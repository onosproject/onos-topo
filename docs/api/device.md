# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/device/device.proto](#api/device/device.proto)
    - [AddRequest](#topo.device.AddRequest)
    - [AddResponse](#topo.device.AddResponse)
    - [Credentials](#topo.device.Credentials)
    - [Device](#topo.device.Device)
    - [Device.AttributesEntry](#topo.device.Device.AttributesEntry)
    - [GetRequest](#topo.device.GetRequest)
    - [GetResponse](#topo.device.GetResponse)
    - [ListRequest](#topo.device.ListRequest)
    - [ListResponse](#topo.device.ListResponse)
    - [ProtocolState](#topo.device.ProtocolState)
    - [RemoveRequest](#topo.device.RemoveRequest)
    - [RemoveResponse](#topo.device.RemoveResponse)
    - [TlsConfig](#topo.device.TlsConfig)
    - [UpdateRequest](#topo.device.UpdateRequest)
    - [UpdateResponse](#topo.device.UpdateResponse)
  
    - [ChannelState](#topo.device.ChannelState)
    - [ConnectivityState](#topo.device.ConnectivityState)
    - [ListResponse.Type](#topo.device.ListResponse.Type)
    - [Protocol](#topo.device.Protocol)
    - [ServiceState](#topo.device.ServiceState)
  
    - [DeviceService](#topo.device.DeviceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/device/device.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/device/device.proto



<a name="topo.device.AddRequest"></a>

### AddRequest
AddRequest adds a device to the topology


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the device to add |






<a name="topo.device.AddResponse"></a>

### AddResponse
AddResponse is sent in response to an AddDeviceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the device with a revision number |






<a name="topo.device.Credentials"></a>

### Credentials
Credentials is the device credentials


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [string](#string) |  | user is the user with which to connect to the device |
| password | [string](#string) |  | password is the password for connecting to the device |






<a name="topo.device.Device"></a>

### Device
Device contains information about a device


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is a globally unique device identifier |
| revision | [uint64](#uint64) |  | revision is the revision of the device |
| address | [string](#string) |  | address is the host:port of the device |
| target | [string](#string) |  | target is the device target |
| version | [string](#string) |  | version is the device software version |
| timeout | [google.protobuf.Duration](#google.protobuf.Duration) |  | timeout indicates the device request timeout |
| credentials | [Credentials](#topo.device.Credentials) |  | credentials contains the credentials for connecting to the device |
| tls | [TlsConfig](#topo.device.TlsConfig) |  | tls is the device TLS configuration |
| type | [string](#string) |  | type is the type of the device |
| role | [string](#string) |  | role is a role for the device |
| attributes | [Device.AttributesEntry](#topo.device.Device.AttributesEntry) | repeated | attributes is an arbitrary mapping of attribute keys/values |
| protocols | [ProtocolState](#topo.device.ProtocolState) | repeated |  |
| displayname | [string](#string) |  | displayname is a user friendly tag |






<a name="topo.device.Device.AttributesEntry"></a>

### Device.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="topo.device.GetRequest"></a>

### GetRequest
GetRequest gets a device by ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is the unique device ID with which to lookup the device |






<a name="topo.device.GetResponse"></a>

### GetResponse
GetResponse carries a device


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the device object |






<a name="topo.device.ListRequest"></a>

### ListRequest
ListRequest requests a stream of devices and changes
By default, the request requests a stream of all devices that are present in the topology when
the request is received by the service. However, if `subscribe` is `true`, the stream will remain
open after all devices have been sent and events that occur following the last device will be
streamed to the client until the stream is closed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscribe | [bool](#bool) |  | subscribe indicates whether to subscribe to events (e.g. ADD, UPDATE, and REMOVE) that occur after all devices have been streamed to the client |






<a name="topo.device.ListResponse"></a>

### ListResponse
ListResponse carries a single device event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ListResponse.Type](#topo.device.ListResponse.Type) |  | type is the type of the event |
| device | [Device](#topo.device.Device) |  | device is the device on which the event occurred |






<a name="topo.device.ProtocolState"></a>

### ProtocolState
ProtocolState contains information related to service and connectivity to a device


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| protocol | [Protocol](#topo.device.Protocol) |  | The protocol to which state relates |
| connectivityState | [ConnectivityState](#topo.device.ConnectivityState) |  | ConnectivityState contains the L3 connectivity information |
| channelState | [ChannelState](#topo.device.ChannelState) |  | ChannelState relates to the availability of the gRPC channel |
| serviceState | [ServiceState](#topo.device.ServiceState) |  | ServiceState indicates the availability of the gRPC servic on top of the channel |






<a name="topo.device.RemoveRequest"></a>

### RemoveRequest
RemoveRequest removes a device by ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the device to remove |






<a name="topo.device.RemoveResponse"></a>

### RemoveResponse
RemoveResponse is sent in response to a RemoveDeviceRequest






<a name="topo.device.TlsConfig"></a>

### TlsConfig
Device TLS configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| caCert | [string](#string) |  | caCert is the name of the device&#39;s CA certificate |
| cert | [string](#string) |  | cert is the name of the device&#39;s certificate |
| key | [string](#string) |  | key is the name of the device&#39;s TLS key |
| plain | [bool](#bool) |  | plain indicates whether to connect to the device over plaintext |
| insecure | [bool](#bool) |  | insecure indicates whether to connect to the device with insecure communication |






<a name="topo.device.UpdateRequest"></a>

### UpdateRequest
UpdateRequest updates a device


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the updated device |






<a name="topo.device.UpdateResponse"></a>

### UpdateResponse
UpdateResponse is sent in response to an UpdateDeviceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device | [Device](#topo.device.Device) |  | device is the device with updated revision |





 


<a name="topo.device.ChannelState"></a>

### ChannelState
ConnectivityState represents the state of a gRPC channel to the device from the service container

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_CHANNEL_STATE | 0 | UNKNOWN_CHANNEL_STATE constant needed to go around proto3 nullifying the 0 values |
| CONNECTED | 1 | CONNECTED indicates the corresponding grpc channel is connected on this device |
| DISCONNECTED | 2 | DISCONNECTED indicates the corresponding grpc channel is not connected on this device |



<a name="topo.device.ConnectivityState"></a>

### ConnectivityState
ConnectivityState represents the L3 reachability of a device from the service container (e.g. enos-config), independently of gRPC or the service itself (e.g. gNMI)

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_CONNECTIVITY_STATE | 0 | UNKNOWN_CONNECTIVITY_STATE constant needed to go around proto3 nullifying the 0 values |
| REACHABLE | 1 | REACHABLE indicates the the service can reach the device at L3 |
| UNREACHABLE | 2 | UNREACHABLE indicates the the service can&#39;t reach the device at L3 |



<a name="topo.device.ListResponse.Type"></a>

### ListResponse.Type
Device event type

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 | NONE indicates this response does not represent a state change |
| ADDED | 1 | ADDED is an event which occurs when a device is added to the topology |
| UPDATED | 2 | UPDATED is an event which occurs when a device is updated |
| REMOVED | 3 | REMOVED is an event which occurs when a device is removed from the topology |



<a name="topo.device.Protocol"></a>

### Protocol
Protocol to interact with a device

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_PROTOCOL | 0 | UNKNOWN_PROTOCOL constant needed to go around proto3 nullifying the 0 values |
| GNMI | 1 | GNMI protocol reference |
| P4RUNTIME | 2 | P4RUNTIME protocol reference |
| GNOI | 3 | GNOI protocol reference |
| E2AP | 4 | E2 Control Plane Protocol |



<a name="topo.device.ServiceState"></a>

### ServiceState
ServiceState represents the state of the gRPC service (e.g. gNMI) to the device from the service container

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_SERVICE_STATE | 0 | UNKNOWN_SERVICE_STATE constant needed to go around proto3 nullifying the 0 values |
| AVAILABLE | 1 | AVAILABLE indicates the corresponding grpc service is available |
| UNAVAILABLE | 2 | UNAVAILABLE indicates the corresponding grpc service is not available |
| CONNECTING | 3 | CONNECTING indicates the corresponding protocol is in the connecting phase on this device |


 

 


<a name="topo.device.DeviceService"></a>

### DeviceService
DeviceService provides an API for managing devices.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Add | [AddRequest](#topo.device.AddRequest) | [AddResponse](#topo.device.AddResponse) | Add adds a device to the topology |
| Update | [UpdateRequest](#topo.device.UpdateRequest) | [UpdateResponse](#topo.device.UpdateResponse) | Update updates a device |
| Get | [GetRequest](#topo.device.GetRequest) | [GetResponse](#topo.device.GetResponse) | Get gets a device by ID |
| List | [ListRequest](#topo.device.ListRequest) | [ListResponse](#topo.device.ListResponse) stream | List gets a stream of device add/update/remove events |
| Remove | [RemoveRequest](#topo.device.RemoveRequest) | [RemoveResponse](#topo.device.RemoveResponse) | Remove removes a device from the topology |

 



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

