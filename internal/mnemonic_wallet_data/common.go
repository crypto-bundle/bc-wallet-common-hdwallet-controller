package mnemonic_wallet_data

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	"github.com/google/uuid"
)

type configurationService interface {
	GetStageName() string
}

type dbStoreService interface {
	AddNewMnemonicWallet(ctx context.Context,
		wallet *entities.MnemonicWallet,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context, uuid uuid.UUID) (*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type cacheStoreService interface {
	SetMnemonicWalletItem(ctx context.Context,
		walletItem *entities.MnemonicWallet,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletItemByUUID(ctx context.Context,
		MnemonicWalletUUID uuid.UUID,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletItemByHash(ctx context.Context,
		MnemonicWalletHash string,
	) (*entities.MnemonicWallet, error)
}
