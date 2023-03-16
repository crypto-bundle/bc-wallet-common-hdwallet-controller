package wallet_data

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common/pkg/postgres"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_data/pgstore"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	persistentStore dbStoreService

	pgConn *postgres.Connection
}

func (s *Service) AddNewWallet(ctx context.Context, walletItem *entities.Wallet) (*entities.Wallet, error) {
	return s.persistentStore.AddNewWallet(ctx, walletItem)
}

func (s *Service) GetWalletByUUID(ctx context.Context, walletUUID string) (*entities.Wallet, error) {
	return s.persistentStore.GetWalletByUUID(ctx, walletUUID)
}

func (s *Service) GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error) {
	return s.persistentStore.GetAllEnabledWallets(ctx)
}

func (s *Service) GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error) {
	return s.persistentStore.GetAllEnabledWalletUUIDList(ctx)
}

func NewService(logger *zap.Logger,
	pgConn *postgres.Connection,
) *Service {
	l := logger.Named("wallet_data.service")
	persistentStoreSrv := pgstore.NewPostgresStore(logger, pgConn)

	return &Service{
		logger:          l,
		persistentStore: persistentStoreSrv,
		pgConn:          pgConn,
	}
}
