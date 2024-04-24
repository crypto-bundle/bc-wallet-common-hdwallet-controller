package config

import "errors"

var (
	ErrWrongMnemonicWordsCount = errors.New("wrong mnemonic words count.allowed: 15, 18, 21, 24")
	ErrBaseAppConfigIsEmpty    = errors.New("base app config is missing or empty")
	ErrMinWordsCount           = errors.New("minimal words count for environment is 21")
)

type MnemonicConfig struct {
	MnemonicWordCount uint8 `envconfig:"HDWALLET_WORDS_COUNT" default:"24"`

	baseAppCfgSrv baseConfigService
}

func (c *MnemonicConfig) GetDefaultMnemonicWordsCount() uint8 {
	return c.MnemonicWordCount
}

//nolint:funlen // its ok
func (c *MnemonicConfig) Prepare() error {
	if c.baseAppCfgSrv == nil {
		return ErrBaseAppConfigIsEmpty
	}

	if c.baseAppCfgSrv.IsProd() && c.MnemonicWordCount <= 18 {
		return ErrMinWordsCount
	}

	switch c.MnemonicWordCount {
	case 15, 18, 21, 24:
		return nil
	default:
		return ErrWrongMnemonicWordsCount
	}
}

//nolint:funlen // its ok
func (c *MnemonicConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
