package wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/google/uuid"
)

type dbStoreService interface {
	AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	UpdateIsEnabledWalletByUUID(ctx context.Context, uuid string, isEnabled bool) error
	GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error)
	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.Wallet, error)
}
