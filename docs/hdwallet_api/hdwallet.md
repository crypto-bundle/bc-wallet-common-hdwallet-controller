# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api.proto](#api.proto)
    - [AddNewWalletRequest](#hdwallet_api.AddNewWalletRequest)
    - [AddNewWalletResponse](#hdwallet_api.AddNewWalletResponse)
    - [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest)
    - [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse)
    - [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity)
    - [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest)
    - [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse)
    - [GetEnabledWalletsRequest](#hdwallet_api.GetEnabledWalletsRequest)
    - [GetEnabledWalletsResponse](#hdwallet_api.GetEnabledWalletsResponse)
  
    - [HdWalletApi](#hdwallet_api.HdWalletApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto



<a name="hdwallet_api.AddNewWalletRequest"></a>

### AddNewWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Title | [string](#string) |  |  |
| Purpose | [string](#string) |  |  |
| IsHot | [bool](#bool) |  |  |






<a name="hdwallet_api.AddNewWalletResponse"></a>

### AddNewWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletUUID | [string](#string) |  |  |






<a name="hdwallet_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccountIndex | [uint32](#uint32) |  |  |
| InternalIndex | [uint32](#uint32) |  |  |
| AddressIndexFrom | [uint32](#uint32) |  |  |
| AddressIndexTo | [uint32](#uint32) |  |  |






<a name="hdwallet_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AddressIdentities | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) | repeated |  |






<a name="hdwallet_api.DerivationAddressIdentity"></a>

### DerivationAddressIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccountIndex | [uint32](#uint32) |  |  |
| InternalIndex | [uint32](#uint32) |  |  |
| AddressIndex | [uint32](#uint32) |  |  |
| Address | [string](#string) |  |  |






<a name="hdwallet_api.DerivationAddressRequest"></a>

### DerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AddressIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.DerivationAddressResponse"></a>

### DerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AddressIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.GetEnabledWalletsRequest"></a>

### GetEnabledWalletsRequest







<a name="hdwallet_api.GetEnabledWalletsResponse"></a>

### GetEnabledWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsUUID | [string](#string) | repeated |  |





 

 

 


<a name="hdwallet_api.HdWalletApi"></a>

### HdWalletApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddNewWallet | [AddNewWalletRequest](#hdwallet_api.AddNewWalletRequest) | [AddNewWalletResponse](#hdwallet_api.AddNewWalletResponse) |  |
| GetEnabledWallets | [GetEnabledWalletsRequest](#hdwallet_api.GetEnabledWalletsRequest) | [GetEnabledWalletsResponse](#hdwallet_api.GetEnabledWalletsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest) | [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse) |  |

 



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

