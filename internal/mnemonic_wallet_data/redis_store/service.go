package redis_store

import (
	"context"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/google/uuid"
	"strings"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	RedisMnemonicWalletPrefix = "mnemonic-wallets"
)

type redisStore struct {
	redisClient *redis.Client
	logger      *zap.Logger

	keyPrefix string
}

func (s *redisStore) SetMnemonicWalletItem(ctx context.Context,
	walletItem *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	rawJSON, err := walletItem.MarshalJSON()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s.%s", s.keyPrefix, walletItem.UUID.String())
	err = s.setMnemonicWalletItemByKey(ctx, key, rawJSON)
	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("%s.%s", s.keyPrefix, walletItem.MnemonicHash)
	err = s.setMnemonicWalletItemByKey(ctx, key, rawJSON)
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func (s *redisStore) setMnemonicWalletItemByKey(ctx context.Context,
	key string,
	rawData []byte,
) error {
	redisCMD := s.redisClient.Set(ctx, key, rawData, 0)
	_, err := redisCMD.Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) GetMnemonicWalletItemByUUID(ctx context.Context,
	MnemonicWalletUUID uuid.UUID,
) (*entities.MnemonicWallet, error) {
	key := fmt.Sprintf("%s.%s", s.keyPrefix, MnemonicWalletUUID.String())

	existCMD := s.redisClient.Exists(ctx, key)
	res, err := existCMD.Result()
	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, nil
	}

	redisCMD := s.redisClient.Get(ctx, key)
	rawJSON, err := redisCMD.Result()
	if err != nil {
		return nil, err
	}

	walletItem := &entities.MnemonicWallet{}
	err = walletItem.UnmarshalJSON([]byte(rawJSON))
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func (s *redisStore) GetMnemonicWalletItemByHash(ctx context.Context,
	MnemonicWalletHash string,
) (*entities.MnemonicWallet, error) {
	key := fmt.Sprintf("%s.%s", s.keyPrefix, MnemonicWalletHash)

	existCMD := s.redisClient.Exists(ctx, key)
	res, err := existCMD.Result()
	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, nil
	}

	redisCMD := s.redisClient.Get(ctx, key)
	rawJSON, err := redisCMD.Result()
	if err != nil {
		return nil, err
	}

	walletItem := &entities.MnemonicWallet{}
	err = walletItem.UnmarshalJSON([]byte(rawJSON))
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func NewRedisStore(logger *zap.Logger,
	cfgSvc configurationService,
	redisClient *redis.Client,
) (*redisStore, error) {
	prefixName := strings.ToUpper(fmt.Sprintf("%s__%s__%s", cfgSvc.GetStageName(), app.ApplicationManagerName,
		RedisMnemonicWalletPrefix),
	)

	return &redisStore{
		redisClient: redisClient,
		logger:      logger,
		keyPrefix:   prefixName,
	}, nil
}
