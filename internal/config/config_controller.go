/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonJWT "github.com/crypto-bundle/bc-wallet-common-lib-jwt/pkg/jwt"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonProfiler "github.com/crypto-bundle/bc-wallet-common-lib-profiler/pkg/profiler"
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
	*commonProfiler.ProfilerConfig
	*commonJWT.JWTConfig
	*hdwallet.HdWalletClientConfig
	*ProcessionEnvironmentConfig
	// -------------------
	// Internal configs
	// -------------------
	EventChannelWorkersCount  int           `envconfig:"EVENT_CHANNEL_WORKERS_COUNT" default:"4"`
	EventChannelBufferSize    int           `envconfig:"EVENT_CHANNEL_BUFFER_SIZE" default:"12"`
	DefaultWalletSessionDelay time.Duration `envconfig:"DEFAULT_WALLET_SESSION_DELAY" default:"2s"`
	DefaultUnloadInterval     time.Duration `envconfig:"DEFAULT_WALLET_UNLOAD_INTERVAL" default:"24s"`
	// VaultCommonTransitKey - common vault transit key for whole processing cluster,
	// must be saved in common vault kv bucket, for example: crypto-bundle/bc-wallet-common/transit
	VaultCommonTransitKey string `envconfig:"VAULT_COMMON_TRANSIT_KEY" secret:"true"`
	// VaultApplicationEncryptionKey - vault encryption key for hd-wallet-controller and hd-wallet-api application,
	// must be saved in bc-wallet-<blockchain_name>-hdwallet vault kv bucket,
	// for example: crypto-bundle/bc-wallet-tron-hdwallet/common
	VaultApplicationEncryptionKey string `envconfig:"VAULT_APP_ENCRYPTION_KEY" secret:"true"`
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

	appName := fmt.Sprintf(ApplicationManagerNameTpl, c.ProcessionEnvironmentConfig.GetNetworkName())

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
