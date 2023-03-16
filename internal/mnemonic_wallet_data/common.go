package mnemonic_wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
)

type dbStoreService interface {
	AddNewMnemonicWallet(ctx context.Context,
		wallet *entities.MnemonicWallet,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}
