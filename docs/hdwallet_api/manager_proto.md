# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [manager_api.proto](#manager_api.proto)
    - [AddNewWalletRequest](#manager_api.AddNewWalletRequest)
    - [AddNewWalletResponse](#manager_api.AddNewWalletResponse)
    - [DerivationAddressByRangeRequest](#manager_api.DerivationAddressByRangeRequest)
    - [DerivationAddressByRangeResponse](#manager_api.DerivationAddressByRangeResponse)
    - [DerivationAddressRequest](#manager_api.DerivationAddressRequest)
    - [DerivationAddressResponse](#manager_api.DerivationAddressResponse)
    - [DisableWalletRequest](#manager_api.DisableWalletRequest)
    - [DisableWalletResponse](#manager_api.DisableWalletResponse)
    - [DisableWalletsRequest](#manager_api.DisableWalletsRequest)
    - [DisableWalletsResponse](#manager_api.DisableWalletsResponse)
    - [EnableWalletRequest](#manager_api.EnableWalletRequest)
    - [EnableWalletResponse](#manager_api.EnableWalletResponse)
    - [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest)
    - [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse)
    - [GetEnabledWalletsResponse.BookmarksEntry](#manager_api.GetEnabledWalletsResponse.BookmarksEntry)
    - [GetWalletInfoRequest](#manager_api.GetWalletInfoRequest)
    - [GetWalletInfoResponse](#manager_api.GetWalletInfoResponse)
    - [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest)
    - [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse)
    - [GetWalletSessionsRequest](#manager_api.GetWalletSessionsRequest)
    - [GetWalletSessionsResponse](#manager_api.GetWalletSessionsResponse)
    - [SignTransactionRequest](#manager_api.SignTransactionRequest)
    - [SignTransactionResponse](#manager_api.SignTransactionResponse)
    - [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest)
    - [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse)
    - [WalletSessionIdentity](#manager_api.WalletSessionIdentity)
  
    - [HdWalletManagerApi](#manager_api.HdWalletManagerApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="manager_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## manager_api.proto



<a name="manager_api.AddNewWalletRequest"></a>

### AddNewWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Title | [string](#string) |  |  |
| Purpose | [string](#string) |  |  |






<a name="manager_api.AddNewWalletResponse"></a>

### AddNewWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| Ranges | [common.RangeRequestUnit](#common.RangeRequestUnit) | repeated |  |






<a name="manager_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentitiesCount | [uint64](#uint64) |  |  |
| AddressIdentities | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) | repeated |  |






<a name="manager_api.DerivationAddressRequest"></a>

### DerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.DerivationAddressResponse"></a>

### DerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |






<a name="manager_api.DisableWalletRequest"></a>

### DisableWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.DisableWalletResponse"></a>

### DisableWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.DisableWalletsRequest"></a>

### DisableWalletsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentities | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="manager_api.DisableWalletsResponse"></a>

### DisableWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentities | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="manager_api.EnableWalletRequest"></a>

### EnableWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.EnableWalletResponse"></a>

### EnableWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetEnabledWalletsRequest"></a>

### GetEnabledWalletsRequest







<a name="manager_api.GetEnabledWalletsResponse"></a>

### GetEnabledWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsCount | [uint32](#uint32) |  |  |
| WalletIdentities | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |
| Bookmarks | [GetEnabledWalletsResponse.BookmarksEntry](#manager_api.GetEnabledWalletsResponse.BookmarksEntry) | repeated |  |






<a name="manager_api.GetEnabledWalletsResponse.BookmarksEntry"></a>

### GetEnabledWalletsResponse.BookmarksEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [uint32](#uint32) |  |  |






<a name="manager_api.GetWalletInfoRequest"></a>

### GetWalletInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletInfoResponse"></a>

### GetWalletInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletSessionRequest"></a>

### GetWalletSessionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletSessionResponse"></a>

### GetWalletSessionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStartedAt | [uint64](#uint64) |  |  |
| SessionExpiredAt | [uint64](#uint64) |  |  |






<a name="manager_api.GetWalletSessionsRequest"></a>

### GetWalletSessionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletSessionsResponse"></a>

### GetWalletSessionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStartedAt | [uint64](#uint64) |  |  |
| SessionExpiredAt | [uint64](#uint64) |  |  |






<a name="manager_api.SignTransactionRequest"></a>

### SignTransactionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletUUID | [string](#string) |  |  |
| MnemonicWalletUUID | [string](#string) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| CreatedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.SignTransactionResponse"></a>

### SignTransactionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SignedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.StartWalletSessionRequest"></a>

### StartWalletSessionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.StartWalletSessionResponse"></a>

### StartWalletSessionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStartedAt | [uint64](#uint64) |  |  |
| SessionExpiredAt | [uint64](#uint64) |  |  |






<a name="manager_api.WalletSessionIdentity"></a>

### WalletSessionIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SessionUUID | [string](#string) |  |  |





 

 

 


<a name="manager_api.HdWalletManagerApi"></a>

### HdWalletManagerApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddNewWallet | [AddNewWalletRequest](#manager_api.AddNewWalletRequest) | [AddNewWalletResponse](#manager_api.AddNewWalletResponse) |  |
| EnableWallet | [EnableWalletRequest](#manager_api.EnableWalletRequest) | [EnableWalletResponse](#manager_api.EnableWalletResponse) |  |
| DisableWallet | [DisableWalletRequest](#manager_api.DisableWalletRequest) | [DisableWalletResponse](#manager_api.DisableWalletResponse) |  |
| DisableWallets | [DisableWalletsRequest](#manager_api.DisableWalletsRequest) | [DisableWalletsResponse](#manager_api.DisableWalletsResponse) |  |
| StartWalletSession | [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest) | [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse) |  |
| GetWalletSession | [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest) | [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse) |  |
| GetAllWalletSessions | [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest) | [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse) |  |
| GetWalletInfo | [GetWalletInfoRequest](#manager_api.GetWalletInfoRequest) | [GetWalletInfoResponse](#manager_api.GetWalletInfoResponse) |  |
| GetEnabledWallets | [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest) | [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#manager_api.DerivationAddressRequest) | [DerivationAddressResponse](#manager_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#manager_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#manager_api.DerivationAddressByRangeResponse) |  |
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

