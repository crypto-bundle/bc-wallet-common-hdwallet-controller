# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [hdwallet_api.proto](#hdwallet_api.proto)
    - [DerivationAddressByRangeRequest](#manager_api.DerivationAddressByRangeRequest)
    - [DerivationAddressByRangeResponse](#manager_api.DerivationAddressByRangeResponse)
    - [DerivationAddressRequest](#manager_api.DerivationAddressRequest)
    - [DerivationAddressResponse](#manager_api.DerivationAddressResponse)
    - [EncryptMnemonicRequest](#manager_api.EncryptMnemonicRequest)
    - [EncryptMnemonicResponse](#manager_api.EncryptMnemonicResponse)
    - [GenerateMnemonicRequest](#manager_api.GenerateMnemonicRequest)
    - [GenerateMnemonicResponse](#manager_api.GenerateMnemonicResponse)
    - [LoadDerivationAddressRequest](#manager_api.LoadDerivationAddressRequest)
    - [LoadDerivationAddressResponse](#manager_api.LoadDerivationAddressResponse)
    - [LoadMnemonicRequest](#manager_api.LoadMnemonicRequest)
    - [LoadMnemonicResponse](#manager_api.LoadMnemonicResponse)
    - [SignTransactionRequest](#manager_api.SignTransactionRequest)
    - [SignTransactionResponse](#manager_api.SignTransactionResponse)
    - [UnLoadMnemonicRequest](#manager_api.UnLoadMnemonicRequest)
    - [UnLoadMnemonicResponse](#manager_api.UnLoadMnemonicResponse)
    - [UnLoadMultipleMnemonicsRequest](#manager_api.UnLoadMultipleMnemonicsRequest)
    - [UnLoadMultipleMnemonicsResponse](#manager_api.UnLoadMultipleMnemonicsResponse)
  
    - [HdWalletApi](#manager_api.HdWalletApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="hdwallet_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## hdwallet_api.proto



<a name="manager_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| Ranges | [common.RangeRequestUnit](#common.RangeRequestUnit) | repeated |  |






<a name="manager_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentitiesCount | [uint64](#uint64) |  |  |
| AddressIdentities | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) | repeated |  |






<a name="manager_api.DerivationAddressRequest"></a>

### DerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.DerivationAddressResponse"></a>

### DerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.EncryptMnemonicRequest"></a>

### EncryptMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| MnemonicData | [bytes](#bytes) |  |  |






<a name="manager_api.EncryptMnemonicResponse"></a>

### EncryptMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="manager_api.GenerateMnemonicRequest"></a>

### GenerateMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GenerateMnemonicResponse"></a>

### GenerateMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="manager_api.LoadDerivationAddressRequest"></a>

### LoadDerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentifier | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.LoadDerivationAddressResponse"></a>

### LoadDerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.LoadMnemonicRequest"></a>

### LoadMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TimeToLive | [uint64](#uint64) |  |  |
| EncryptedMnemonicData | [bytes](#bytes) |  |  |






<a name="manager_api.LoadMnemonicResponse"></a>

### LoadMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.SignTransactionRequest"></a>

### SignTransactionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentifier | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| CreatedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.SignTransactionResponse"></a>

### SignTransactionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicWalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SignedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.UnLoadMnemonicRequest"></a>

### UnLoadMnemonicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.UnLoadMnemonicResponse"></a>

### UnLoadMnemonicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.UnLoadMultipleMnemonicsRequest"></a>

### UnLoadMultipleMnemonicsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="manager_api.UnLoadMultipleMnemonicsResponse"></a>

### UnLoadMultipleMnemonicsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |





 

 

 


<a name="manager_api.HdWalletApi"></a>

### HdWalletApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GenerateMnemonic | [GenerateMnemonicRequest](#manager_api.GenerateMnemonicRequest) | [GenerateMnemonicResponse](#manager_api.GenerateMnemonicResponse) |  |
| EncryptMnemonic | [EncryptMnemonicRequest](#manager_api.EncryptMnemonicRequest) | [EncryptMnemonicResponse](#manager_api.EncryptMnemonicResponse) |  |
| LoadMnemonic | [LoadMnemonicRequest](#manager_api.LoadMnemonicRequest) | [LoadMnemonicResponse](#manager_api.LoadMnemonicResponse) |  |
| UnLoadMnemonic | [UnLoadMnemonicRequest](#manager_api.UnLoadMnemonicRequest) | [UnLoadMnemonicResponse](#manager_api.UnLoadMnemonicResponse) |  |
| UnLoadMultipleMnemonics | [UnLoadMultipleMnemonicsRequest](#manager_api.UnLoadMultipleMnemonicsRequest) | [UnLoadMultipleMnemonicsResponse](#manager_api.UnLoadMultipleMnemonicsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#manager_api.DerivationAddressRequest) | [DerivationAddressResponse](#manager_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#manager_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#manager_api.DerivationAddressByRangeResponse) |  |
| LoadDerivationAddress | [LoadDerivationAddressRequest](#manager_api.LoadDerivationAddressRequest) | [LoadDerivationAddressResponse](#manager_api.LoadDerivationAddressResponse) |  |
| SignTransaction | [SignTransactionRequest](#manager_api.SignTransactionRequest) | [SignTransactionResponse](#manager_api.SignTransactionResponse) |  |

 



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

