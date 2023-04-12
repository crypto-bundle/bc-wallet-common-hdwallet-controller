package wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_data/pgstore"
	"github.com/google/uuid"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	persistentStore dbStoreService

	pgConn *commonPostgres.Connection
}

func (s *Service) AddNewWallet(ctx context.Context, walletItem *entities.Wallet) (*entities.Wallet, error) {
	return s.persistentStore.AddNewWallet(ctx, walletItem)
}

func (s *Service) SetEnabledToWalletByUUID(ctx context.Context, uuid string) error {
	return s.persistentStore.UpdateIsEnabledWalletByUUID(ctx, uuid, true)
}

func (s *Service) GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.Wallet, error) {
	return s.persistentStore.GetWalletByUUID(ctx, walletUUID)
}

func (s *Service) GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error) {
	return s.persistentStore.GetAllEnabledWallets(ctx)
}

func (s *Service) GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error) {
	return s.persistentStore.GetAllEnabledWalletUUIDList(ctx)
}

func NewService(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *Service {
	l := logger.Named("wallet_data.service")
	persistentStoreSrv := pgstore.NewPostgresStore(logger, pgConn)

	return &Service{
		logger:          l,
		persistentStore: persistentStoreSrv,
		pgConn:          pgConn,
	}
}
