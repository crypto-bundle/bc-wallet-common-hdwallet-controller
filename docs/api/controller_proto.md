# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [manager_api.proto](#manager_api.proto)
    - [AddNewWalletRequest](#manager_api.AddNewWalletRequest)
    - [AddNewWalletResponse](#manager_api.AddNewWalletResponse)
    - [CloseWalletSessionsRequest](#manager_api.CloseWalletSessionsRequest)
    - [CloseWalletSessionsResponse](#manager_api.CloseWalletSessionsResponse)
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
    - [ExecuteSignRequestReq](#manager_api.ExecuteSignRequestReq)
    - [ExecuteSignRequestResponse](#manager_api.ExecuteSignRequestResponse)
    - [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest)
    - [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse)
    - [GetEnabledWalletsResponse.BookmarksEntry](#manager_api.GetEnabledWalletsResponse.BookmarksEntry)
    - [GetWalletInfoRequest](#manager_api.GetWalletInfoRequest)
    - [GetWalletInfoResponse](#manager_api.GetWalletInfoResponse)
    - [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest)
    - [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse)
    - [GetWalletSessionsRequest](#manager_api.GetWalletSessionsRequest)
    - [GetWalletSessionsResponse](#manager_api.GetWalletSessionsResponse)
    - [ImportWalletRequest](#manager_api.ImportWalletRequest)
    - [ImportWalletResponse](#manager_api.ImportWalletResponse)
    - [PrepareSignRequestReq](#manager_api.PrepareSignRequestReq)
    - [PrepareSignRequestResponse](#manager_api.PrepareSignRequestResponse)
    - [SessionInfo](#manager_api.SessionInfo)
    - [SignPurposeIdentity](#manager_api.SignPurposeIdentity)
    - [SignRequestData](#manager_api.SignRequestData)
    - [SignRequestIdentity](#manager_api.SignRequestIdentity)
    - [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest)
    - [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse)
    - [WalletSessionIdentity](#manager_api.WalletSessionIdentity)
  
    - [SignRequestData.ReqStatus](#manager_api.SignRequestData.ReqStatus)
  
    - [HdWalletManagerApi](#manager_api.HdWalletManagerApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="manager_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## manager_api.proto



<a name="manager_api.AddNewWalletRequest"></a>

### AddNewWalletRequest







<a name="manager_api.AddNewWalletResponse"></a>

### AddNewWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.CloseWalletSessionsRequest"></a>

### CloseWalletSessionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |






<a name="manager_api.CloseWalletSessionsResponse"></a>

### CloseWalletSessionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |






<a name="manager_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| Ranges | [common.RangeRequestUnit](#common.RangeRequestUnit) | repeated |  |






<a name="manager_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
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






<a name="manager_api.ExecuteSignRequestReq"></a>

### ExecuteSignRequestReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SignRequestIdentifier | [SignRequestIdentity](#manager_api.SignRequestIdentity) |  |  |
| CreatedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.ExecuteSignRequestResponse"></a>

### ExecuteSignRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SignRequestIdentifier | [SignRequestIdentity](#manager_api.SignRequestIdentity) |  |  |
| SignedTxData | [bytes](#bytes) |  |  |






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
| Session | [SessionInfo](#manager_api.SessionInfo) |  |  |






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
| ActiveSessions | [SessionInfo](#manager_api.SessionInfo) | repeated |  |






<a name="manager_api.ImportWalletRequest"></a>

### ImportWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicPhrase | [bytes](#bytes) |  |  |






<a name="manager_api.ImportWalletResponse"></a>

### ImportWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.PrepareSignRequestReq"></a>

### PrepareSignRequestReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SignPurposeIdentifier | [SignPurposeIdentity](#manager_api.SignPurposeIdentity) |  |  |






<a name="manager_api.PrepareSignRequestResponse"></a>

### PrepareSignRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicIdentity | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| TxOwnerIdentity | [common.DerivationAddressIdentity](#common.DerivationAddressIdentity) |  |  |
| SignatureRequestInfo | [SignRequestData](#manager_api.SignRequestData) |  |  |






<a name="manager_api.SessionInfo"></a>

### SessionInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SessionIdentity | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStartedAt | [uint64](#uint64) |  |  |
| SessionExpiredAt | [uint64](#uint64) |  |  |






<a name="manager_api.SignPurposeIdentity"></a>

### SignPurposeIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UUID | [string](#string) |  |  |






<a name="manager_api.SignRequestData"></a>

### SignRequestData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Identifier | [SignRequestIdentity](#manager_api.SignRequestIdentity) |  |  |
| Status | [SignRequestData.ReqStatus](#manager_api.SignRequestData.ReqStatus) |  |  |
| CreateAt | [uint64](#uint64) |  |  |






<a name="manager_api.SignRequestIdentity"></a>

### SignRequestIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UUID | [string](#string) |  |  |






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





 


<a name="manager_api.SignRequestData.ReqStatus"></a>

### SignRequestData.ReqStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| REQUEST_PLACEHOLDER | 0 |  |
| REQUEST_CREATED | 1 |  |
| REQUEST_PREPARED | 2 |  |
| REQUEST_SIGNED | 3 |  |
| REQUEST_FAILED | 4 |  |


 

 


<a name="manager_api.HdWalletManagerApi"></a>

### HdWalletManagerApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddNewWallet | [AddNewWalletRequest](#manager_api.AddNewWalletRequest) | [AddNewWalletResponse](#manager_api.AddNewWalletResponse) |  |
| ImportWallet | [ImportWalletRequest](#manager_api.ImportWalletRequest) | [ImportWalletResponse](#manager_api.ImportWalletResponse) |  |
| EnableWallet | [EnableWalletRequest](#manager_api.EnableWalletRequest) | [EnableWalletResponse](#manager_api.EnableWalletResponse) |  |
| GetWalletInfo | [GetWalletInfoRequest](#manager_api.GetWalletInfoRequest) | [GetWalletInfoResponse](#manager_api.GetWalletInfoResponse) |  |
| GetEnabledWallets | [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest) | [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse) |  |
| DisableWallet | [DisableWalletRequest](#manager_api.DisableWalletRequest) | [DisableWalletResponse](#manager_api.DisableWalletResponse) |  |
| DisableWallets | [DisableWalletsRequest](#manager_api.DisableWalletsRequest) | [DisableWalletsResponse](#manager_api.DisableWalletsResponse) |  |
| StartWalletSession | [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest) | [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse) |  |
| GetWalletSession | [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest) | [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse) |  |
| GetAllWalletSessions | [GetWalletSessionsRequest](#manager_api.GetWalletSessionsRequest) | [GetWalletSessionsResponse](#manager_api.GetWalletSessionsResponse) |  |
| CloseWalletSession | [CloseWalletSessionsRequest](#manager_api.CloseWalletSessionsRequest) | [CloseWalletSessionsResponse](#manager_api.CloseWalletSessionsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#manager_api.DerivationAddressRequest) | [DerivationAddressResponse](#manager_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#manager_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#manager_api.DerivationAddressByRangeResponse) |  |
| PrepareSignRequest | [PrepareSignRequestReq](#manager_api.PrepareSignRequestReq) | [PrepareSignRequestResponse](#manager_api.PrepareSignRequestResponse) |  |
| ExecuteSignRequest | [ExecuteSignRequestReq](#manager_api.ExecuteSignRequestReq) | [ExecuteSignRequestResponse](#manager_api.ExecuteSignRequestResponse) |  |

 



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

