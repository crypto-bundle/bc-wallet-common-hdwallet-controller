package redis_store

import (
	"context"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"strings"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	RedisMnemonicWalletPrefix         = "mnemonic-wallets"
	RedisMnemonicWalletSessionsPrefix = "mnemonic-wallets-sessions"
)

type redisStore struct {
	redisClient *redis.Client
	logger      *zap.Logger

	walletInfoKeyPrefix     string
	walletSessionsKeyPrefix string
}

func (s *redisStore) SetMnemonicWalletItem(ctx context.Context,
	walletItem *entities.MnemonicWallet,
) error {
	rawJSON, err := walletItem.MarshalJSON()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, walletItem.UUID.String())
	err = s.setMnemonicWalletItemByKey(ctx, key, rawJSON)
	if err != nil {
		return err
	}

	return nil
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

func (s *redisStore) GetAllWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	walletSearchKey := fmt.Sprintf("prefix:%s.*", s.walletInfoKeyPrefix)
	iter := s.redisClient.Scan(ctx, 0, walletSearchKey, 0).Iterator()

	walletList := make([]*entities.MnemonicWallet, 0)
	for iter.Next(ctx) {
		item := &entities.MnemonicWallet{}

		loopErr := item.UnmarshalJSON([]byte(iter.Val()))
		if loopErr != nil {
			return nil, loopErr
		}

		walletList = append(walletList, item)
	}
	err := iter.Err()
	if err != nil {
		return nil, err
	}

	if len(walletList) == 0 {
		return nil, nil
	}

	return walletList, nil
}

func (s *redisStore) GetMnemonicWalletByUUID(ctx context.Context,
	MnemonicWalletUUID string,
) (*entities.MnemonicWallet, error) {
	walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, MnemonicWalletUUID)
	rawJSON, err := s.redisClient.Get(ctx, walletKey).Result()
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

func (s *redisStore) GetMnemonicWalletInfoByUUID(ctx context.Context,
	MnemonicWalletUUID string,
) (wallet *entities.MnemonicWallet, sessions []*entities.MnemonicWalletSession, err error) {
	resList, err := s.redisClient.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, MnemonicWalletUUID)
		walletSessionsKey := fmt.Sprintf("prefix:%s.%s-*", s.walletSessionsKeyPrefix,
			MnemonicWalletUUID)

		pipeliner.Get(ctx, walletKey)
		pipeliner.Scan(ctx, 0, walletSessionsKey, 0)

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	err = resList[0].Err()
	if err != nil {
		return nil, nil, err
	}

	err = resList[1].Err()
	if err != nil {
		return nil, nil, err
	}

	rawJSON := resList[0].(*redis.StringCmd).Val()
	err = wallet.UnmarshalJSON([]byte(rawJSON))
	if err != nil {
		return nil, nil, err
	}

	iter := resList[1].(*redis.ScanCmd).Iterator()
	if err != nil {
		return nil, nil, err
	}

	sessionList := make([]*entities.MnemonicWalletSession, 0)
	for iter.Next(ctx) {
		session := &entities.MnemonicWalletSession{}

		loopErr := session.UnmarshalJSON([]byte(iter.Val()))
		if loopErr != nil {
			return nil, nil, loopErr
		}

		sessions = append(sessions, session)
	}
	err = iter.Err()
	if err != nil {
		return nil, nil, err
	}

	if len(sessionList) > 0 {
		sessions = sessionList
	}

	return wallet, sessions, nil
}

func (s *redisStore) GetMnemonicWalletSessionInfoByUUID(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (wallet *entities.MnemonicWallet, session *entities.MnemonicWalletSession, err error) {
	resList, err := s.redisClient.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, walletUUID)
		walletSessionsKey := fmt.Sprintf("%s.%s", s.walletSessionsKeyPrefix, sessionUUID)

		pipeliner.Get(ctx, walletKey)
		pipeliner.Get(ctx, walletSessionsKey)

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	err = resList[0].Err()
	if err != nil {
		return nil, nil, err
	}

	err = resList[1].Err()
	if err != nil {
		return nil, nil, err
	}

	rawJSON := resList[0].(*redis.StringCmd).Val()
	err = wallet.UnmarshalJSON([]byte(rawJSON))
	if err != nil {
		return nil, nil, err
	}

	rawJSON = resList[1].(*redis.StringCmd).Val()
	err = session.UnmarshalJSON([]byte(rawJSON))
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *redisStore) SetMnemonicWalletSessionItem(ctx context.Context,
	sessionItem *entities.MnemonicWalletSession,
) error {
	rawJSON, err := sessionItem.MarshalJSON()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s.%s-%s", s.walletSessionsKeyPrefix,
		sessionItem.MnemonicWalletUUID, sessionItem.UUID)
	cmd := s.redisClient.Set(ctx, key, rawJSON, sessionItem.ExpiredAt.Sub(time.Now()))

	_, err = cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) UnsetWalletSession(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) error {
	key := fmt.Sprintf("%s.%s-%s", s.walletSessionsKeyPrefix,
		walletUUID, sessionUUID)
	cmd := s.redisClient.Del(ctx, key)

	_, err := cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) FullUnsetMnemonicWallet(ctx context.Context,
	mnemonicWalletUUID string,
) error {
	walletKey := fmt.Sprintf("%s.%s", s.walletSessionsKeyPrefix, mnemonicWalletUUID)

	walletSessionsSearchKeyPattern := fmt.Sprintf("prefix:%s.%s-*", s.walletSessionsKeyPrefix,
		mnemonicWalletUUID)
	keysCmd := s.redisClient.Keys(ctx, walletSessionsSearchKeyPattern)
	err := keysCmd.Err()
	if err != nil {
		return err
	}

	keysForDelete := keysCmd.Val()
	keysForDelete = append(keysForDelete, walletKey)

	delCmd := s.redisClient.Del(ctx, keysForDelete...)
	err = delCmd.Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) UnsetMultipleWallets(ctx context.Context,
	mnemonicWalletsUUIDs []string,
	sessionsUUIDs []string,
) error {
	var j int

	deleteKeys := make([]string, len(mnemonicWalletsUUIDs)+len(sessionsUUIDs))
	for i, _ := range mnemonicWalletsUUIDs {
		deleteKeys[j] = fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, mnemonicWalletsUUIDs[i])
		j++
	}

	for i, _ := range sessionsUUIDs {
		deleteKeys[j] = fmt.Sprintf("%s.%s", s.walletSessionsKeyPrefix, sessionsUUIDs[i])
		j++
	}

	_, err := s.redisClient.Del(ctx, deleteKeys...).Result()
	if err != nil {
		return err
	}

	return nil
}

func NewRedisStore(logger *zap.Logger,
	cfgSvc configurationService,
	redisClient *redis.Client,
) *redisStore {
	prefixName := strings.ToUpper(fmt.Sprintf("%s__%s__%s", cfgSvc.GetStageName(),
		app.ApplicationManagerName,
		RedisMnemonicWalletPrefix),
	)

	sessionsPrefixName := strings.ToUpper(fmt.Sprintf("%s__%s__%s", cfgSvc.GetStageName(),
		app.ApplicationManagerName,
		RedisMnemonicWalletSessionsPrefix),
	)

	return &redisStore{
		redisClient:             redisClient,
		logger:                  logger,
		walletInfoKeyPrefix:     prefixName,
		walletSessionsKeyPrefix: sessionsPrefixName,
	}
}
