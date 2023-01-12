package config

import (
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"

	"github.com/crypto-bundle/bc-wallet-common/pkg/crypter"

	"github.com/kelseyhightower/envconfig"
)

type MnemonicConfig struct {
	Mnemonic   string
	Hash       string
	BlockChain string
}

func (c *MnemonicConfig) GetHash() string {
	return c.Hash
}

type HDWalletConfig struct {
	crypter.Config
}

//nolint:funlen // its ok
func (c *HDWalletConfig) Prepare() error {
	err := envconfig.Process(app.HDWalletConfigPrefix, c)
	if err != nil {
		return err
	}

	return c.Config.Prepare()
}
