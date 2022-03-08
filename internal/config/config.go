package config

import (
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"
	"github.com/kelseyhightower/envconfig"
)

// Config for application
type Config struct {
	// -------------------
	// Application configs
	// -------------------
	*BaseConfig
	// -------------------
	// Database config
	// -------------------
	*DBConfig
	// -------------------
	// GRPC service config
	// -------------------
	*GrpcConfig
	// -------------------
	// HD wallet config
	// -------------------
	*HDWalletConfig

	VaultClient vaulter
}

// Prepare variables to static configuration
func (c *Config) Prepare() error {
	err := envconfig.Process(app.ApplicationName, c)
	if err != nil {
		return err
	}

	c.BaseConfig = &BaseConfig{}
	err = c.BaseConfig.Prepare()
	if err != nil {
		return err
	}

	//b, err := c.VaultClient.GetCredentialsBytes()
	//if err != nil {
	//	return err
	//}

	c.DBConfig = &DBConfig{
		vaultData: nil, // temporary unless vault is not working
	}
	err = c.DBConfig.Prepare()
	if err != nil {
		return err
	}

	c.HDWalletConfig = &HDWalletConfig{}
	err = c.HDWalletConfig.Prepare()
	if err != nil {
		return err
	}

	c.GrpcConfig = &GrpcConfig{}
	err = c.GrpcConfig.Prepare()
	if err != nil {
		return err
	}

	return err
}
