# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [controller_api.proto](#controller_api.proto)
    - [AccessTokenData](#manager_api.AccessTokenData)
    - [AccessTokenIdentity](#manager_api.AccessTokenIdentity)
    - [AddNewWalletRequest](#manager_api.AddNewWalletRequest)
    - [AddNewWalletResponse](#manager_api.AddNewWalletResponse)
    - [AppInstanceIdentity](#manager_api.AppInstanceIdentity)
    - [CloseWalletSessionsRequest](#manager_api.CloseWalletSessionsRequest)
    - [CloseWalletSessionsResponse](#manager_api.CloseWalletSessionsResponse)
    - [DisableWalletRequest](#manager_api.DisableWalletRequest)
    - [DisableWalletResponse](#manager_api.DisableWalletResponse)
    - [DisableWalletsRequest](#manager_api.DisableWalletsRequest)
    - [DisableWalletsResponse](#manager_api.DisableWalletsResponse)
    - [DisableWalletsResponse.BookmarksEntry](#manager_api.DisableWalletsResponse.BookmarksEntry)
    - [EnableWalletRequest](#manager_api.EnableWalletRequest)
    - [EnableWalletResponse](#manager_api.EnableWalletResponse)
    - [EnableWalletsRequest](#manager_api.EnableWalletsRequest)
    - [EnableWalletsResponse](#manager_api.EnableWalletsResponse)
    - [EnableWalletsResponse.BookmarksEntry](#manager_api.EnableWalletsResponse.BookmarksEntry)
    - [Event](#manager_api.Event)
    - [ExecuteSignRequestReq](#manager_api.ExecuteSignRequestReq)
    - [ExecuteSignRequestResponse](#manager_api.ExecuteSignRequestResponse)
    - [GetAccountRequest](#manager_api.GetAccountRequest)
    - [GetAccountResponse](#manager_api.GetAccountResponse)
    - [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest)
    - [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse)
    - [GetEnabledWalletsResponse.BookmarksEntry](#manager_api.GetEnabledWalletsResponse.BookmarksEntry)
    - [GetMultipleAccountRequest](#manager_api.GetMultipleAccountRequest)
    - [GetMultipleAccountResponse](#manager_api.GetMultipleAccountResponse)
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
    - [SignRequestEvent](#manager_api.SignRequestEvent)
    - [SignRequestIdentity](#manager_api.SignRequestIdentity)
    - [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest)
    - [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse)
    - [WalletSessionEvent](#manager_api.WalletSessionEvent)
    - [WalletSessionIdentity](#manager_api.WalletSessionIdentity)
  
    - [Event.Type](#manager_api.Event.Type)
    - [SignRequestData.ReqStatus](#manager_api.SignRequestData.ReqStatus)
    - [SignRequestEvent.Type](#manager_api.SignRequestEvent.Type)
    - [WalletSessionEvent.Type](#manager_api.WalletSessionEvent.Type)
    - [WalletSessionStatus](#manager_api.WalletSessionStatus)
  
    - [HdWalletControllerApi](#manager_api.HdWalletControllerApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="controller_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## controller_api.proto



<a name="manager_api.AccessTokenData"></a>

### AccessTokenData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccessTokenIdentifier | [AccessTokenIdentity](#manager_api.AccessTokenIdentity) |  |  |
| AccessTokenData | [bytes](#bytes) |  |  |






<a name="manager_api.AccessTokenIdentity"></a>

### AccessTokenIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UUID | [string](#string) |  |  |






<a name="manager_api.AddNewWalletRequest"></a>

### AddNewWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccessTokens | [AccessTokenData](#manager_api.AccessTokenData) | repeated |  |






<a name="manager_api.AddNewWalletResponse"></a>

### AddNewWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| WalletStatus | [common.WalletStatus](#common.WalletStatus) |  |  |






<a name="manager_api.AppInstanceIdentity"></a>

### AppInstanceIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UUID | [string](#string) |  |  |






<a name="manager_api.CloseWalletSessionsRequest"></a>

### CloseWalletSessionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |






<a name="manager_api.CloseWalletSessionsResponse"></a>

### CloseWalletSessionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStatus | [WalletSessionStatus](#manager_api.WalletSessionStatus) |  |  |






<a name="manager_api.DisableWalletRequest"></a>

### DisableWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.DisableWalletResponse"></a>

### DisableWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| WalletStatus | [common.WalletStatus](#common.WalletStatus) |  |  |






<a name="manager_api.DisableWalletsRequest"></a>

### DisableWalletsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifiers | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="manager_api.DisableWalletsResponse"></a>

### DisableWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsCount | [uint32](#uint32) |  |  |
| WalletsData | [common.MnemonicWalletData](#common.MnemonicWalletData) | repeated |  |
| Bookmarks | [DisableWalletsResponse.BookmarksEntry](#manager_api.DisableWalletsResponse.BookmarksEntry) | repeated |  |






<a name="manager_api.DisableWalletsResponse.BookmarksEntry"></a>

### DisableWalletsResponse.BookmarksEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [uint32](#uint32) |  |  |






<a name="manager_api.EnableWalletRequest"></a>

### EnableWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.EnableWalletResponse"></a>

### EnableWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| WalletStatus | [common.WalletStatus](#common.WalletStatus) |  |  |






<a name="manager_api.EnableWalletsRequest"></a>

### EnableWalletsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifiers | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) | repeated |  |






<a name="manager_api.EnableWalletsResponse"></a>

### EnableWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsCount | [uint32](#uint32) |  |  |
| WalletsData | [common.MnemonicWalletData](#common.MnemonicWalletData) | repeated |  |
| Bookmarks | [EnableWalletsResponse.BookmarksEntry](#manager_api.EnableWalletsResponse.BookmarksEntry) | repeated |  |






<a name="manager_api.EnableWalletsResponse.BookmarksEntry"></a>

### EnableWalletsResponse.BookmarksEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [uint32](#uint32) |  |  |






<a name="manager_api.Event"></a>

### Event



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| EventType | [Event.Type](#manager_api.Event.Type) |  |  |
| AppInstanceIdentifier | [AppInstanceIdentity](#manager_api.AppInstanceIdentity) |  |  |
| Data | [bytes](#bytes) |  |  |






<a name="manager_api.ExecuteSignRequestReq"></a>

### ExecuteSignRequestReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SignRequestIdentifier | [SignRequestIdentity](#manager_api.SignRequestIdentity) |  |  |
| CreatedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.ExecuteSignRequestResponse"></a>

### ExecuteSignRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |
| SignatureRequestInfo | [SignRequestData](#manager_api.SignRequestData) |  |  |
| SignedTxData | [bytes](#bytes) |  |  |






<a name="manager_api.GetAccountRequest"></a>

### GetAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |






<a name="manager_api.GetAccountResponse"></a>

### GetAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |






<a name="manager_api.GetEnabledWalletsRequest"></a>

### GetEnabledWalletsRequest







<a name="manager_api.GetEnabledWalletsResponse"></a>

### GetEnabledWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsCount | [uint32](#uint32) |  |  |
| WalletsData | [common.MnemonicWalletData](#common.MnemonicWalletData) | repeated |  |
| Bookmarks | [GetEnabledWalletsResponse.BookmarksEntry](#manager_api.GetEnabledWalletsResponse.BookmarksEntry) | repeated |  |






<a name="manager_api.GetEnabledWalletsResponse.BookmarksEntry"></a>

### GetEnabledWalletsResponse.BookmarksEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [uint32](#uint32) |  |  |






<a name="manager_api.GetMultipleAccountRequest"></a>

### GetMultipleAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| Parameters | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="manager_api.GetMultipleAccountResponse"></a>

### GetMultipleAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AccountIdentitiesCount | [uint64](#uint64) |  |  |
| AccountIdentifiers | [common.AccountIdentity](#common.AccountIdentity) | repeated |  |






<a name="manager_api.GetWalletInfoRequest"></a>

### GetWalletInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletInfoResponse"></a>

### GetWalletInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| WalletStatus | [common.WalletStatus](#common.WalletStatus) |  |  |






<a name="manager_api.GetWalletSessionRequest"></a>

### GetWalletSessionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |






<a name="manager_api.GetWalletSessionResponse"></a>

### GetWalletSessionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| Session | [SessionInfo](#manager_api.SessionInfo) |  |  |






<a name="manager_api.GetWalletSessionsRequest"></a>

### GetWalletSessionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.GetWalletSessionsResponse"></a>

### GetWalletSessionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| ActiveSessions | [SessionInfo](#manager_api.SessionInfo) | repeated |  |






<a name="manager_api.ImportWalletRequest"></a>

### ImportWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| MnemonicPhrase | [bytes](#bytes) |  |  |
| AccessTokens | [AccessTokenData](#manager_api.AccessTokenData) | repeated |  |






<a name="manager_api.ImportWalletResponse"></a>

### ImportWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.PrepareSignRequestReq"></a>

### PrepareSignRequestReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SignPurposeIdentifier | [SignPurposeIdentity](#manager_api.SignPurposeIdentity) |  |  |






<a name="manager_api.PrepareSignRequestResponse"></a>

### PrepareSignRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| AccountIdentifier | [common.AccountIdentity](#common.AccountIdentity) |  |  |
| SignatureRequestInfo | [SignRequestData](#manager_api.SignRequestData) |  |  |






<a name="manager_api.SessionInfo"></a>

### SessionInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStatus | [WalletSessionStatus](#manager_api.WalletSessionStatus) |  |  |
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






<a name="manager_api.SignRequestEvent"></a>

### SignRequestEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| EventType | [SignRequestEvent.Type](#manager_api.SignRequestEvent.Type) |  |  |
| SignRequestIdentifier | [SignRequestIdentity](#manager_api.SignRequestIdentity) |  |  |






<a name="manager_api.SignRequestIdentity"></a>

### SignRequestIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UUID | [string](#string) |  |  |






<a name="manager_api.StartWalletSessionRequest"></a>

### StartWalletSessionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |






<a name="manager_api.StartWalletSessionResponse"></a>

### StartWalletSessionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |
| SessionStatus | [WalletSessionStatus](#manager_api.WalletSessionStatus) |  |  |
| SessionStartedAt | [uint64](#uint64) |  |  |
| SessionExpiredAt | [uint64](#uint64) |  |  |






<a name="manager_api.WalletSessionEvent"></a>

### WalletSessionEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| EventType | [WalletSessionEvent.Type](#manager_api.WalletSessionEvent.Type) |  |  |
| WalletIdentifier | [common.MnemonicWalletIdentity](#common.MnemonicWalletIdentity) |  |  |
| SessionIdentifier | [WalletSessionIdentity](#manager_api.WalletSessionIdentity) |  |  |






<a name="manager_api.WalletSessionIdentity"></a>

### WalletSessionIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SessionUUID | [string](#string) |  |  |





 


<a name="manager_api.Event.Type"></a>

### Event.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| EVENT_TYPE_PLACEHOLDER | 0 |  |
| EVENT_TYPE_SESSION | 1 |  |
| EVENT_TYPE_SIGN_REQUEST | 2 |  |



<a name="manager_api.SignRequestData.ReqStatus"></a>

### SignRequestData.ReqStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| REQUEST_PLACEHOLDER | 0 |  |
| REQUEST_CREATED | 1 |  |
| REQUEST_PREPARED | 2 |  |
| REQUEST_SIGNED | 3 |  |
| REQUEST_FAILED | 4 |  |



<a name="manager_api.SignRequestEvent.Type"></a>

### SignRequestEvent.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| PLACEHOLDER | 0 |  |
| PREPARED | 1 |  |
| CLOSED | 2 |  |



<a name="manager_api.WalletSessionEvent.Type"></a>

### WalletSessionEvent.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| PLACEHOLDER | 0 |  |
| STARTED | 1 |  |
| CLOSED | 2 |  |



<a name="manager_api.WalletSessionStatus"></a>

### WalletSessionStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| WALLET_SESSION_STATUS_PLACEHOLDER | 0 |  |
| WALLET_SESSION_STATUS_PREPARED | 1 |  |
| WALLET_SESSION_STATUS_CLOSED | 2 |  |


 

 


<a name="manager_api.HdWalletControllerApi"></a>

### HdWalletControllerApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddNewWallet | [AddNewWalletRequest](#manager_api.AddNewWalletRequest) | [AddNewWalletResponse](#manager_api.AddNewWalletResponse) |  |
| ImportWallet | [ImportWalletRequest](#manager_api.ImportWalletRequest) | [ImportWalletResponse](#manager_api.ImportWalletResponse) |  |
| EnableWallet | [EnableWalletRequest](#manager_api.EnableWalletRequest) | [EnableWalletResponse](#manager_api.EnableWalletResponse) |  |
| GetWalletInfo | [GetWalletInfoRequest](#manager_api.GetWalletInfoRequest) | [GetWalletInfoResponse](#manager_api.GetWalletInfoResponse) |  |
| GetEnabledWallets | [GetEnabledWalletsRequest](#manager_api.GetEnabledWalletsRequest) | [GetEnabledWalletsResponse](#manager_api.GetEnabledWalletsResponse) |  |
| DisableWallet | [DisableWalletRequest](#manager_api.DisableWalletRequest) | [DisableWalletResponse](#manager_api.DisableWalletResponse) |  |
| DisableWallets | [DisableWalletsRequest](#manager_api.DisableWalletsRequest) | [DisableWalletsResponse](#manager_api.DisableWalletsResponse) |  |
| EnableWallets | [EnableWalletsRequest](#manager_api.EnableWalletsRequest) | [EnableWalletsResponse](#manager_api.EnableWalletsResponse) |  |
| StartWalletSession | [StartWalletSessionRequest](#manager_api.StartWalletSessionRequest) | [StartWalletSessionResponse](#manager_api.StartWalletSessionResponse) |  |
| GetWalletSession | [GetWalletSessionRequest](#manager_api.GetWalletSessionRequest) | [GetWalletSessionResponse](#manager_api.GetWalletSessionResponse) |  |
| GetAllWalletSessions | [GetWalletSessionsRequest](#manager_api.GetWalletSessionsRequest) | [GetWalletSessionsResponse](#manager_api.GetWalletSessionsResponse) |  |
| CloseWalletSession | [CloseWalletSessionsRequest](#manager_api.CloseWalletSessionsRequest) | [CloseWalletSessionsResponse](#manager_api.CloseWalletSessionsResponse) |  |
| GetAccount | [GetAccountRequest](#manager_api.GetAccountRequest) | [GetAccountResponse](#manager_api.GetAccountResponse) |  |
| GetMultipleAccounts | [GetMultipleAccountRequest](#manager_api.GetMultipleAccountRequest) | [GetMultipleAccountResponse](#manager_api.GetMultipleAccountResponse) |  |
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

