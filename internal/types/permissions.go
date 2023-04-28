package types

type PermissionType uint8

const (
	PermissionTypeAddNewWalletName           = "permission_add_new_wallet"
	PermissionTypeGetWalletName              = "permission_get_wallet"
	PermissionTypeGetEnabledWalletsName      = "permission_get_enabled_wallets"
	PermissionTypeGetDerivationAddressName   = "permission_get_derivation_address"
	PermissionTypeGetDerivationAddressesName = "permission_get_derivation_addresses"
	PermissionTypeStartWalletSessionName     = "permission_start_wallet_session"
	PermissionTypeGetWalletSessionName       = "permission_get_wallet_session"
	PermissionTypeSignTransactionName        = "permission_sign_transaction"
)

const (
	PermissionTypePlaceHolder PermissionType = iota
	PermissionTypeAddNewWallet
	PermissionTypeGetWallet
	PermissionTypeGetEnabledWallets
	PermissionTypeGetDerivationAddress
	PermissionTypeGetDerivationAddresses
	PermissionTypeStartWalletSession
	PermissionTypeGetWalletSession
	PermissionTypeSignTransaction
)

func (d PermissionType) String() string {
	switch d {
	case PermissionTypeAddNewWallet:
		return PermissionTypeAddNewWalletName
	case PermissionTypeGetWallet:
		return PermissionTypeGetWalletName
	case PermissionTypeGetEnabledWallets:
		return PermissionTypeGetEnabledWalletsName
	case PermissionTypeGetDerivationAddress:
		return PermissionTypeGetDerivationAddressName
	case PermissionTypeGetDerivationAddresses:
		return PermissionTypeGetDerivationAddressesName
	case PermissionTypeStartWalletSession:
		return PermissionTypeStartWalletSessionName
	case PermissionTypeGetWalletSession:
		return PermissionTypeGetWalletSessionName
	case PermissionTypeSignTransaction:
		return PermissionTypeSignTransactionName
	case PermissionTypePlaceHolder:
		fallthrough
	default:
		return ""
	}
}
