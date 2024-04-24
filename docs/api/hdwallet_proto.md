# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [hdwallet_api.proto](#hdwallet_api.proto)
    - [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest)
    - [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse)
    - [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest)
    - [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse)
    - [EncryptMnemonicRequest](#hdwallet_api.EncryptMnemonicRequest)
    - [EncryptMnemonicResponse](#hdwallet_api.EncryptMnemonicResponse)
    - [GenerateMnemonicRequest](#hdwallet_api.GenerateMnemonicRequest)
    - [GenerateMnemonicResponse](#hdwallet_api.GenerateMnemonicResponse)
    - [LoadDerivationAddressRequest](#hdwallet_api.LoadDerivationAddressRequest)
    - [LoadDerivationAddressResponse](#hdwallet_api.LoadDerivationAddressResponse)
    - [LoadMnemonicRequest](#hdwallet_api.LoadMnemonicRequest)
    - [LoadMnemonicResponse](#hdwallet_api.LoadMnemonicResponse)
    - [SignDataRequest](#hdwallet_api.SignDataRequest)
    - [SignDataResponse](#hdwallet_api.SignDataResponse)
    - [UnLoadMnemonicRequest](#hdwallet_api.UnLoadMnemonicRequest)
    - [UnLoadMnemonicResponse](#hdwallet_api.UnLoadMnemonicResponse)
    - [UnLoadMultipleMnemonicsRequest](#hdwallet_api.UnLoadMultipleMnemonicsRequest)
    - [UnLoadMultipleMnemonicsResponse](#hdwallet_api.UnLoadMultipleMnemonicsResponse)
    - [ValidateMnemonicRequest](#hdwallet_api.ValidateMnemonicRequest)
    - [ValidateMnemonicResponse](#hdwallet_api.ValidateMnemonicResponse)
  
    - [HdWalletApi](#hdwallet_api.HdWalletApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="hdwallet_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## hdwallet_api.proto



<a name="hdwallet_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| Ranges | [common.RangeRequestUnit](#common.RangeRequestUnit) | repeated |  |






<a name="hdwallet_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentitiesCount | [uint64](#uint64) |  |  |
| AddressIdentities | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) | repeated |  |






<a name="hdwallet_api.DerivationAddressRequest"></a>

### DerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.DerivationAddressResponse"></a>

### DerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.EncryptMnemonicRequest"></a>

### EncryptMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| MnemonicData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.EncryptMnemonicResponse"></a>

### EncryptMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.GenerateMnemonicRequest"></a>

### GenerateMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="hdwallet_api.GenerateMnemonicResponse"></a>

### GenerateMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.LoadDerivationAddressRequest"></a>

### LoadDerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentifier | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.LoadDerivationAddressResponse"></a>

### LoadDerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.LoadMnemonicRequest"></a>

### LoadMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TimeToLive | [uint64](#uint64) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.LoadMnemonicResponse"></a>

### LoadMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="hdwallet_api.SignDataRequest"></a>

### SignDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentifier | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| DataForSign | [bytes](#bytes) |  |  |






<a name="hdwallet_api.SignDataResponse"></a>

### SignDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SignedData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.UnLoadMnemonicRequest"></a>

### UnLoadMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="hdwallet_api.UnLoadMnemonicResponse"></a>

### UnLoadMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="hdwallet_api.UnLoadMultipleMnemonicsRequest"></a>

### UnLoadMultipleMnemonicsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="hdwallet_api.UnLoadMultipleMnemonicsResponse"></a>

### UnLoadMultipleMnemonicsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="hdwallet_api.ValidateMnemonicRequest"></a>

### ValidateMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| MnemonicData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.ValidateMnemonicResponse"></a>

### ValidateMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| IsValid | [bool](#bool) |  |  |





 

 

 


<a name="hdwallet_api.HdWalletApi"></a>

### HdWalletApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GenerateMnemonic | [GenerateMnemonicRequest](#hdwallet_api.GenerateMnemonicRequest) | [GenerateMnemonicResponse](#hdwallet_api.GenerateMnemonicResponse) |  |
| EncryptMnemonic | [EncryptMnemonicRequest](#hdwallet_api.EncryptMnemonicRequest) | [EncryptMnemonicResponse](#hdwallet_api.EncryptMnemonicResponse) |  |
| ValidateMnemonic | [ValidateMnemonicRequest](#hdwallet_api.ValidateMnemonicRequest) | [ValidateMnemonicResponse](#hdwallet_api.ValidateMnemonicResponse) |  |
| LoadMnemonic | [LoadMnemonicRequest](#hdwallet_api.LoadMnemonicRequest) | [LoadMnemonicResponse](#hdwallet_api.LoadMnemonicResponse) |  |
| UnLoadMnemonic | [UnLoadMnemonicRequest](#hdwallet_api.UnLoadMnemonicRequest) | [UnLoadMnemonicResponse](#hdwallet_api.UnLoadMnemonicResponse) |  |
| UnLoadMultipleMnemonics | [UnLoadMultipleMnemonicsRequest](#hdwallet_api.UnLoadMultipleMnemonicsRequest) | [UnLoadMultipleMnemonicsResponse](#hdwallet_api.UnLoadMultipleMnemonicsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest) | [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse) |  |
| LoadDerivationAddress | [LoadDerivationAddressRequest](#hdwallet_api.LoadDerivationAddressRequest) | [LoadDerivationAddressResponse](#hdwallet_api.LoadDerivationAddressResponse) |  |
| SignData | [SignDataRequest](#hdwallet_api.SignDataRequest) | [SignDataResponse](#hdwallet_api.SignDataResponse) |  |

 



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

