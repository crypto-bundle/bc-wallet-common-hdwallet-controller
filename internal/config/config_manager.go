package config

import (
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
	"strings"
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
	// GRPCBindRaw port string, default "8080"
	GRPCBindRaw string `envconfig:"API_GRPC_PORT" default:"8080"`
	// ----------------------------
	// Calculated config parameters
	GRPCBind string
	// ----------------------------
	// Dependencies
	baseAppCfgSrv baseConfigService
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
