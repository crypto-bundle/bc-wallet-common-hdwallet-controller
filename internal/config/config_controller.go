package config

import (
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"strings"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"

	"github.com/google/uuid"
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
	EventChannelWorkersCount  int           `envconfig:"EVENT_CHANNEL_WORKERS_COUNT" default:"4"`
	EventChannelBufferSize    int           `envconfig:"EVENT_CHANNEL_BUFFER_SIZE" default:"12"`
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
	GRPCBind         string
	InstanceUUID     uuid.UUID
	EventChannelName string
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

func (c *MangerConfig) GetInstanceIdentifier() uuid.UUID {
	return c.InstanceUUID
}

func (c *MangerConfig) GetEventChannelName() string {
	return c.EventChannelName
}
func (c *MangerConfig) GetEventChannelWorkersCount() int {
	return c.EventChannelWorkersCount
}
func (c *MangerConfig) GetEventChannelBufferSize() int {
	return c.EventChannelBufferSize
}

// Prepare variables to static configuration
func (c *MangerConfig) Prepare() error {
	c.GRPCBind = fmt.Sprintf(":%s", strings.TrimLeft(c.GRPCBindRaw, ":"))
	c.InstanceUUID = uuid.New()

	appName := fmt.Sprintf(app.ApplicationManagerNameTpl, c.ProcessionEnvironmentConfig.GetNetworkName())

	c.baseAppCfgSrv.SetApplicationName(appName)

	c.EventChannelName = strings.ToUpper(fmt.Sprintf("%s__%s__%s__%s", c.GetStageName(),
		c.GetApplicationName(), c.GetProviderName(), c.GetNetworkName()))

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
