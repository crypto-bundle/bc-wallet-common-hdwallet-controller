/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
