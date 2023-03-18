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

package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
)

type configService interface {
	GetDefaultHotWalletUnloadInterval() time.Duration
	GetDefaultWalletUnloadInterval() time.Duration

	GetMnemonicsCountPerWallet() uint8
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type walletsDataService interface {
	AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	GetWalletByUUID(ctx context.Context, uuid string) (*entities.Wallet, error)
	GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error)
}

type walletPoolUnitMakerService interface {
	CreateWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title, purpose string,
	) (WalletPoolUnitService, error)
}

type WalletPoolUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	AddMnemonicUnit(unit walletPoolMnemonicUnitService) error
	GetWalletUUID() uuid.UUID
	GetWalletTitle() string
	GetWalletPurpose() string
	GetWalletPublicData() *types.PublicWalletData
	GetAddressByPath(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		accountIndex uint32,
		internalIndex uint32,
		addressIndexFrom uint32,
		addressIndexTo uint32,
	) ([]*types.PublicDerivationAddressData, error)
	SignTransaction(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transaction *tronCore.Transaction,
	) (*types.PublicSignTxData, error)
}

type walletPoolMnemonicUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	LoadWallet(ctx context.Context) error
	UnloadWallet(ctx context.Context) error
	GetPublicData() *types.PublicMnemonicWalletData

	IsHotWalletUnit() bool
	GetMnemonicUUID() uuid.UUID
	GetAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		accountIndex uint32,
		internalIndex uint32,
		addressIndexFrom uint32,
		addressIndexTo uint32,
	) ([]*types.PublicDerivationAddressData, error)
	SignTransaction(ctx context.Context,
		account, change, index uint32,
		transaction *tronCore.Transaction,
	) (*types.PublicSignTxData, error)
}

type walletPoolInitService interface {
	LoadAndInitWallets(ctx context.Context) error
	GetWalletPoolUnits() map[uuid.UUID]WalletPoolUnitService
}

type walletPoolService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	SetWalletUnits(ctx context.Context,
		walletUnits map[uuid.UUID]WalletPoolUnitService,
	) error
	AddAndStartWalletUnit(ctx context.Context,
		walletUUID uuid.UUID,
		walletUnit WalletPoolUnitService,
	) error
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
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

type mnemonicWalletConfig interface {
	GetMnemonicWalletPurpose() string
	GetMnemonicWalletHash() string
	IsHotWallet() bool
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}

type walleter interface {
	GetAddress() (string, error)
	GetPubKey() string
	GetPrvKey() (string, error)
	GetPath() string
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewTronWallet(account, change, address uint32) (*hdwallet.Tron, error)
}

type mnemonicGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type encryptService interface {
	Encrypt(msg string) (string, error)
	Decrypt(encMsg string) ([]byte, error)
}
