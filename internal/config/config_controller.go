/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
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
	// System access token
	// must be saved in bc-wallet-common vault kv bucket,
	// for example: kv/data/crypto-bundle/bc-wallet-common/jwt
	SystemAccessTokenHash string `envconfig:"JWT_SYSTEM_ACCESS_TOKEN_HASH" secret:"true"`
	// ManagerApiGRPCBindRaw - address and port string for tcp bind
	ManagerApiGRPCBindRaw string `envconfig:"MANAGER_API_GRPC_PORT" default:"8080"`
	// WalletApiGRPCBindRaw - address and port string for tcp bind
	WalletApiGRPCBindRaw string `envconfig:"WALLET_API_GRPC_PORT" default:"8081"`
	// ----------------------------
	// Calculated config parameters
	ManagerApiGRPCBind string
	WalletApiGRPCBind  string
	InstanceUUID       uuid.UUID
	EventChannelName   string
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

func (c *MangerConfig) GetSystemAccessTokenHash() string {
	return c.SystemAccessTokenHash
}

func (c *MangerConfig) GetVaultAppEncryptionKey() string {
	return c.VaultApplicationEncryptionKey
}

func (c *MangerConfig) GetManagerApiBindAddress() string {
	return c.ManagerApiGRPCBind
}

func (c *MangerConfig) GetWalletApiBindAddress() string {
	return c.WalletApiGRPCBind
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
	c.ManagerApiGRPCBind = fmt.Sprintf(":%s", strings.TrimLeft(c.ManagerApiGRPCBindRaw, ":"))
	c.WalletApiGRPCBind = fmt.Sprintf(":%s", strings.TrimLeft(c.WalletApiGRPCBindRaw, ":"))

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
