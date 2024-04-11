package config

import (
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
	"strings"
	"time"
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
	*hdwallet.HdWalletClientConfig
	*ProcessionEnvironmentConfig
	// -------------------
	// Internal configs
	// -------------------
	DefaultWalletSessionDelay time.Duration `envconfig:"DEFAULT_WALLET_SESSION_DELAY" default:"2s"`
	DefaultUnloadInterval     time.Duration `envconfig:"DEFAULT_WALLET_UNLOAD_INTERVAL" default:"24s"`
	// VaultCommonTransitKey - common vault transit key for whole processing cluster
	VaultCommonTransitKey string `envconfig:"VAULT_COMMON_TRANSIT_KEY" default:"-"`
	// VaultApplicationEncryptionKey - vault encryption key for hd-wallet-controller and hd-wallet-api application
	VaultApplicationEncryptionKey string `envconfig:"VAULT_APP_ENCRYPTION_KEY" default:"-"`
	// GRPCBindRaw port string, default "8080"
	GRPCBindRaw string `envconfig:"API_GRPC_PORT" default:"8080"`
	// ----------------------------
	// Calculated config parameters
	GRPCBind string
	// ----------------------------
	// Dependencies
	baseAppCfgSrv baseConfigService
}

func (c *MangerConfig) GetDefaultWalletSessionDelay() time.Duration {
	return c.DefaultWalletSessionDelay
}

func (c *MangerConfig) GetDefaultWalletUnloadInterval() time.Duration {
	return c.DefaultUnloadInterval
}

func (c *MangerConfig) GetVaultCommonTransit() string {
	return c.VaultCommonTransitKey
}

func (c *MangerConfig) GetVaultAppEncryptionKey() string {
	return c.VaultApplicationEncryptionKey
}

func (c *MangerConfig) GetBindPort() string {
	return c.GRPCBind
}

// Prepare variables to static configuration
func (c *MangerConfig) Prepare() error {
	c.GRPCBind = fmt.Sprintf(":%s", strings.TrimLeft(c.GRPCBindRaw, ":"))

	return nil
}

func (c *MangerConfig) PrepareWith(cfgSvcList ...interface{}) error {
	for _, cfgSrv := range cfgSvcList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
