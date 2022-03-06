package config

import (
	"crypto/sha256"
	"fmt"

	"bc-wallet-eth-hdwallet/internal/app"

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
	ETHMnemonic string `envconfig:"HDWALLET_ETH_MNEMONIC"`

	EthMnemonicConfig *MnemonicConfig
}

//nolint:funlen // its ok
func (c *HDWalletConfig) Prepare() error {
	err := envconfig.Process(app.HDWalletConfigPrefix, c)
	if err != nil {
		return err
	}

	// ETH wallet
	ethWalletMnemoStr := c.GetETHMnemonics()
	ethWalletHash := sha256.Sum256([]byte(ethWalletMnemoStr))
	c.EthMnemonicConfig = &MnemonicConfig{
		Mnemonic:   c.GetETHMnemonics(),
		Hash:       fmt.Sprintf("%x", ethWalletHash),
		BlockChain: app.BlockChainName,
	}

	return nil
}

func (c *HDWalletConfig) GetETHMnemonics() string {
	return c.ETHMnemonic
}

func (c *HDWalletConfig) GetETHMnemonicsConfig() *MnemonicConfig {
	return c.EthMnemonicConfig
}
