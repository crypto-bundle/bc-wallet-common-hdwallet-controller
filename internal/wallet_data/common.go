package wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
)

type dbStoreService interface {
	AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error)
	GetWalletByUUID(ctx context.Context, uuid string) (*entities.Wallet, error)
}
