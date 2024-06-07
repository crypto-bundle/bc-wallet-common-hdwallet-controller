/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package mnemonic_wallet_data

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/nats_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/pg_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/redis_store"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	txStmtSvc transactionalStatementManager

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

func (s *Service) DisableWallet(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletUUIDTag, walletUUID)

	err = s.txStmtSvc.BeginTxWithRollbackOnError(tCtx, func(txStmtCtx context.Context) error {
		updatedItem, clbErr := s.persistentStore.UpdateWalletStatus(txStmtCtx,
			walletUUID, types.MnemonicWalletStatusDisabled)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return err
		}

	}


	go func(item entities.MnemonicWallet) {
		_, err = s.redisCacheStore.SetMnemonicWalletItem(context.Background(), &item)
		if err != nil {
			s.logger.Error("unable to save mnemonic wallet item in redis cache store", zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, item.UUID.String()),
				zap.String(app.WalletUUIDTag, item.WalletUUID.String()))
		}
	}(*updatedItem)

	go func(item entities.MnemonicWallet) {
		_, err = s.natsKVCacheStore.SetMnemonicWalletItem(context.Background(), &item)
		if err != nil {
			s.logger.Error("unable to save mnemonic wallet item in nats kv cache store", zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, item.UUID.String()),
				zap.String(app.WalletUUIDTag, item.WalletUUID.String()))
		}
	}(*updatedItem)

	return updatedItem, nil
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
