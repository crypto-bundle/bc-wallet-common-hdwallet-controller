/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

package handlers

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
)

const (
	MethodNameTag = "method_name"
)

type walleter interface {
	GetAddressByPath(ctx context.Context,
		walletUUID string,
		account, change, index uint32,
	) (string, error)

	GetAddressByPathByRange(ctx context.Context,
		walletUUID string,
		accountIndex uint32,
		internalIndex uint32,
		addressIndexFrom uint32,
		addressIndexTo uint32,
	) ([]*types.PublicDerivationAddressData, error)

	GetEnabledWalletsUUID(ctx context.Context) ([]string, error)

	CreateNewWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title string,
		purpose string,
	) (*types.PublicWalletData, error)
}
