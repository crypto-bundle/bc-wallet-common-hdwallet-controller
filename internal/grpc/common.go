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

package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
)

type walletManagerService interface {
	CreateNewWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title string,
		purpose string,
	) (*types.PublicWalletData, error)
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*types.PublicDerivationAddressData, error)

	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		accountIndex uint32,
		internalIndex uint32,
		addressIndexFrom uint32,
		addressIndexTo uint32,
	) ([]*types.PublicDerivationAddressData, error)

	GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error)

	SignTransaction(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transaction *tronCore.Transaction,
	) (*types.PublicSignTxData, error)
}

type marshallerService interface {
	MarshallCreateWalletData(*types.PublicWalletData) (*pbApi.AddNewWalletResponse, error)
	MarshallGetAddressData(*types.PublicDerivationAddressData) (*pbApi.DerivationAddressResponse, error)
	MarshallGetAddressByRange([]*types.PublicDerivationAddressData) (*pbApi.DerivationAddressByRangeResponse, error)
	MarshallGetEnabledWallets([]*types.PublicWalletData) (*pbApi.GetEnabledWalletsResponse, error)
	MarshallSignTransaction(
		publicSignTxData *types.PublicSignTxData,
	) (*pbApi.SignTransactionResponse, error)
}
