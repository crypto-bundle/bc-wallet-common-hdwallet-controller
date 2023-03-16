package mnemonic_wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common/pkg/postgres"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data/pgstore"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	persistentStore dbStoreService

	pgConn *postgres.Connection
}

func (s *Service) AddNewMnemonicWallet(ctx context.Context,
	wallet *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	return s.persistentStore.AddNewMnemonicWallet(ctx, wallet)
}

func (s *Service) GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error) {
	return s.persistentStore.GetMnemonicWalletByHash(ctx, hash)
}

func (s *Service) GetMnemonicWalletUUID(ctx context.Context, walletUUID string) (*entities.MnemonicWallet, error) {
	return s.persistentStore.GetMnemonicWalletByUUID(ctx, walletUUID)
}

func (s *Service) GetMnemonicWalletsByUUIDList(ctx context.Context,
	UUIDList []string,
) ([]*entities.MnemonicWallet, error) {
	return s.persistentStore.GetMnemonicWalletsByUUIDList(ctx, UUIDList)
}

func (s *Service) GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	return s.persistentStore.GetAllHotMnemonicWallets(ctx)
}

func (s *Service) GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	return s.persistentStore.GetAllNonHotMnemonicWallets(ctx)
}

func NewService(logger *zap.Logger,
	pgConn *postgres.Connection,
) *Service {
	l := logger.Named("mnemonic_wallet_data.service")
	persistentStoreSrv := pgstore.NewPostgresStore(logger, pgConn)

	return &Service{
		logger:          l,
		persistentStore: persistentStoreSrv,
		pgConn:          pgConn,
	}
}
