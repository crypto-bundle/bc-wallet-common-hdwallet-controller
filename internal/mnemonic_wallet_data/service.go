package mnemonic_wallet_data

import (
	"context"

	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data/nats_store"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data/pg_store"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data/redis_store"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	persistentStore  dbStoreService
	redisCacheStore  cacheStoreService
	natsKVCacheStore cacheStoreService
}

func (s *Service) AddNewMnemonicWallet(ctx context.Context,
	wallet *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	mnemoWalletItem, err := s.persistentStore.AddNewMnemonicWallet(ctx, wallet)
	if err != nil {
		s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemoWalletItem.UUID.String()),
			zap.String(app.WalletUUIDTag, mnemoWalletItem.WalletUUID.String()))

		return nil, err
	}

	go func(item entities.MnemonicWallet) {
		_, err = s.redisCacheStore.SetMnemonicWalletItem(context.Background(), &item)
		if err != nil {
			s.logger.Error("unable to save mnemonic wallet item in redis cache store", zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, item.UUID.String()),
				zap.String(app.WalletUUIDTag, item.WalletUUID.String()))
		}
	}(*mnemoWalletItem)

	go func(item entities.MnemonicWallet) {
		_, err = s.natsKVCacheStore.SetMnemonicWalletItem(context.Background(), &item)
		if err != nil {
			s.logger.Error("unable to save mnemonic wallet item in nats kv cache store", zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, item.UUID.String()),
				zap.String(app.WalletUUIDTag, item.WalletUUID.String()))
		}
	}(*mnemoWalletItem)

	return mnemoWalletItem, nil
}

func (s *Service) GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletHashTag, hash)

	walletItem, err := s.redisCacheStore.GetMnemonicWalletItemByHash(ctx, hash)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from redis cache store", zap.Error(err),
			zap.String(app.MnemonicWalletHashTag, hash))
	}

	if walletItem != nil {
		return walletItem, nil
	}

	walletItem, err = s.natsKVCacheStore.GetMnemonicWalletItemByHash(ctx, hash)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from nats cache store", zap.Error(err),
			zap.String(app.MnemonicWalletHashTag, hash))
	}

	if walletItem != nil {
		return walletItem, nil
	}

	walletItem, err = s.persistentStore.GetMnemonicWalletByHash(ctx, hash)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from persistent store", zap.Error(err),
			zap.String(app.MnemonicWalletHashTag, hash))
	}

	if walletItem == nil {
		return nil, nil
	}

	return walletItem, nil
}

func (s *Service) GetMnemonicWalletUUID(ctx context.Context,
	mnemoWalletUUID uuid.UUID,
) (*entities.MnemonicWallet, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletUUIDTag, mnemoWalletUUID.String())

	walletItem, err := s.redisCacheStore.GetMnemonicWalletItemByUUID(ctx, mnemoWalletUUID)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from redis cache store", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemoWalletUUID.String()))
	}

	if walletItem != nil {
		return walletItem, nil
	}

	walletItem, err = s.natsKVCacheStore.GetMnemonicWalletItemByUUID(ctx, mnemoWalletUUID)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from nats cache store", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemoWalletUUID.String()))
	}

	if walletItem != nil {
		return walletItem, nil
	}

	walletItem, err = s.persistentStore.GetMnemonicWalletByUUID(ctx, mnemoWalletUUID)
	if err != nil {
		s.logger.Error("unable get mnemonicWalletItem from persistent store", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemoWalletUUID.String()))
	}

	if walletItem == nil {
		return nil, nil
	}

	return walletItem, nil
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
	configSvc configurationService,
	pgConn *commonPostgres.Connection,
	redisClient *redis.Client,
	natsConn *commonNats.Connection,
) (*Service, error) {
	l := logger.Named("mnemonic_wallet_data.service")
	persistentStoreSrv := pg_store.NewPostgresStore(logger, pgConn)

	redisStore, err := redis_store.NewRedisStore(logger, configSvc, redisClient)
	if err != nil {
		return nil, err
	}

	natsJetSteamContext, err := natsConn.GetConnection().JetStream()
	if err != nil {
		return nil, err
	}

	natsKvStore, err := nats_store.NewNatsStore(logger, configSvc, natsJetSteamContext)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger: l,

		persistentStore:  persistentStoreSrv,
		redisCacheStore:  redisStore,
		natsKVCacheStore: natsKvStore,
	}, nil
}
