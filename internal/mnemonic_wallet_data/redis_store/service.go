package redis_store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	RedisMnemonicWalletPrefix         = "mnemonic-wallets"
	RedisMnemonicWalletSessionsPrefix = "mnemonic-wallets-sessions"
)

var (
	ErrUnableCastRedisValue = errors.New("unable to cast returned redis value")
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
		if errors.Is(redis.Nil, err) {
			return nil, nil
		}

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
	walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, MnemonicWalletUUID)
	walletSessionsKey := fmt.Sprintf("%s.%s.*", s.walletSessionsKeyPrefix,
		MnemonicWalletUUID)

	iterator := s.redisClient.Scan(ctx, 0, walletSessionsKey, 0).Iterator()
	keys := []string{
		walletKey,
	}
	for iterator.Next(ctx) {
		keys = append(keys, iterator.Val())
	}

	resList, err := s.redisClient.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.MGet(ctx, keys...)

		return nil
	})
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil, nil
		}

		return nil, nil, err
	}

	err = resList[0].Err()
	if err != nil {
		return nil, nil, err
	}

	// len(resList)-1 - cuz first N values its redis.ScanCmd
	resultSliceCmd, ok := resList[len(resList)-1].(*redis.SliceCmd)
	if !ok {
		return nil, nil, ErrUnableCastRedisValue
	}
	resultList := resultSliceCmd.Val()
	count := len(resultList)

	switch true {
	case count > 1: // wallet and sessions data found
		sessionList := make([]*entities.MnemonicWalletSession, count-1)
		for i := 1; i != count; i++ {
			value, casted := resultList[i].(string)
			if !casted {
				return nil, nil, ErrUnableCastRedisValue
			}

			session := &entities.MnemonicWalletSession{}
			loopErr := session.UnmarshalJSON([]byte(value))
			if loopErr != nil {
				return nil, nil, loopErr
			}

			sessionList[i-1] = session
		}

		sessions = sessionList

		fallthrough
	case count == 1: // only wallet data found
		if resultList[0] == nil {
			break
		}

		walletData, casted := resultList[0].(string) // wallet item
		if !casted {
			return nil, nil, ErrUnableCastRedisValue
		}
		walletItem := &entities.MnemonicWallet{}
		caseErr := walletItem.UnmarshalJSON([]byte(walletData))
		if caseErr != nil {
			return nil, nil, caseErr
		}

		wallet = walletItem
	case count == 1 && resultList[0] == nil: // empty data
		return nil, nil, nil
	}

	return
}

func (s *redisStore) GetMnemonicWalletSessionInfoByUUID(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (walletItem *entities.MnemonicWallet, sessionItem *entities.MnemonicWalletSession, err error) {
	resList, err := s.redisClient.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, walletUUID)
		walletSessionsKey := fmt.Sprintf("%s.%s.%s", s.walletSessionsKeyPrefix, walletUUID, sessionUUID)

		pipeliner.MGet(ctx, walletKey, walletSessionsKey)

		return nil
	})
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	err = resList[0].Err()
	if err != nil {
		return nil, nil, err
	}

	list := resList[0].(*redis.SliceCmd).Val()

	if list[0] != nil {
		rawJSON, ok := list[0].(string)
		if !ok {
			return nil, nil, ErrUnableCastRedisValue
		}

		wallet := &entities.MnemonicWallet{}
		err = wallet.UnmarshalJSON([]byte(rawJSON))
		if err != nil {
			return nil, nil, err
		}

		walletItem = wallet
	}

	if list[1] != nil {
		rawJSON, ok := list[1].(string)
		if !ok {
			return nil, nil, ErrUnableCastRedisValue
		}

		session := &entities.MnemonicWalletSession{}
		err = session.UnmarshalJSON([]byte(rawJSON))
		if err != nil {
			return nil, nil, err
		}

		sessionItem = session
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

	key := fmt.Sprintf("%s.%s.%s", s.walletSessionsKeyPrefix,
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
	key := fmt.Sprintf("%s.%s.%s", s.walletSessionsKeyPrefix,
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
	walletKey := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, mnemonicWalletUUID)

	walletSessionsSearchKeyPattern := fmt.Sprintf("%s.%s.*", s.walletSessionsKeyPrefix,
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
		cfgSvc.GetApplicationName(),
		RedisMnemonicWalletPrefix),
	)

	sessionsPrefixName := strings.ToUpper(fmt.Sprintf("%s__%s__%s", cfgSvc.GetStageName(),
		cfgSvc.GetApplicationName(),
		RedisMnemonicWalletSessionsPrefix),
	)

	return &redisStore{
		redisClient:             redisClient,
		logger:                  logger,
		walletInfoKeyPrefix:     prefixName,
		walletSessionsKeyPrefix: sessionsPrefixName,
	}
}
