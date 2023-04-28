package config

import (
	"fmt"
	"strings"
	"time"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
)

// MangerConfig for application
type MangerConfig struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonHealthcheck.HealthcheckHTTPConfig
	*VaultWrappedConfig
	*commonPostgres.PostgresConfig
	*commonNats.NatsConfig
	*commonRedis.RedisConfig
	// -------------------
	// Internal configs
	// -------------------
	*GrpcConfig
	// GRPCBindRaw port string, default "8080"
	GRPCBindRaw                         string        `envconfig:"API_GRPC_PORT" default:"8080"`
	WalletManagerUnloadHotInterval      time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_HOT_INTERVAL" default:"15s"`
	WalletManagerUnloadInterval         time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_INTERVAL" default:"8s"`
	WalletManagerMnemonicPerWalletCount uint8         `envconfig:"WALLET_MANAGER_MNEMONICS_PER_WALLET_COUNT" default:"3"`
	// ----------------------------
	// Calculated config parameters
	GRPCBind string
}

func (c *MangerConfig) GetBindPort() string {
	return c.GRPCBind
}

func (c *MangerConfig) GetDefaultHotWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadHotInterval
}

func (c *MangerConfig) GetDefaultWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadInterval
}

func (c *MangerConfig) GetMnemonicsCountPerWallet() uint8 {
	return c.WalletManagerMnemonicPerWalletCount
}

// Prepare variables to static configuration
func (c *MangerConfig) Prepare() error {
	c.GRPCBind = fmt.Sprintf(":%s", strings.TrimLeft(c.GRPCBindRaw, ":"))

	return nil
}

func (c *MangerConfig) PrepareWith(cfgSvcList ...interface{}) error {
	return nil
}
