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

func (s *redisStore) SetMultipleMnemonicWallets(ctx context.Context,
	walletItems []*entities.MnemonicWallet,
) error {
	_, err := s.redisClient.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, walletItem := range walletItems {
			rawJSON, loopErr := walletItem.MarshalJSON()
			if loopErr != nil {
				return loopErr
			}

			key := fmt.Sprintf("%s.%s", s.walletInfoKeyPrefix, walletItem.UUID.String())
			pipeliner.Set(ctx, key, rawJSON, 0)
		}

		return nil
	})
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
	walletSearchKey := fmt.Sprintf("%s.*", s.walletInfoKeyPrefix)
	sliceCmd := s.redisClient.Keys(ctx, walletSearchKey)
	keys := sliceCmd.Val()

	dataListRaw := s.redisClient.MGet(ctx, keys...)
	resList, err := dataListRaw.Result()
	if err != nil {
		return nil, err
	}

	walletList := make([]*entities.MnemonicWallet, len(resList))
	for i := 0; i != len(resList); i++ {
		item := &entities.MnemonicWallet{}

		JSONstr, isCasted := resList[i].(string)
		if !isCasted {
			return nil, ErrUnableCastRedisValue
		}

		loopErr := item.UnmarshalJSON([]byte(JSONstr))
		if loopErr != nil {
			return nil, loopErr
		}

		walletList[i] = item
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
	toUnsetList []string,
) error {
	deleteKeys := make([]string, 0)
	for _, walletUUID := range toUnsetList {
		deleteKeys = append(deleteKeys, fmt.Sprintf("%s.%s",
			s.walletInfoKeyPrefix, walletUUID))
	}

	_, err := s.redisClient.Del(ctx, deleteKeys...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) UnsetMultipleSessions(ctx context.Context,
	toUnsetMap map[string][]string,
) error {
	deleteKeys := make([]string, 0)
	for walletUUID, sessionUUIDList := range toUnsetMap {
		for _, sessionUUID := range sessionUUIDList {
			deleteKeys = append(deleteKeys, fmt.Sprintf("%s.%s.%s",
				s.walletSessionsKeyPrefix, walletUUID, sessionUUID))
		}
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
