package config

import (
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/app"

	"github.com/cryptowize-tech/bc-wallet-common/pkg/postgres"

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
	*postgres.PostgresConfig
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

	c.PostgresConfig = &postgres.PostgresConfig{}
	err = c.PostgresConfig.Prepare()
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
