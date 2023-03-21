/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

package config

import (
	"time"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

// Config for application
type Config struct {
	// -------------------
	// Application configs
	// -------------------
	*commonConfig.BaseConfig
	// -------------------
	// Logger configs
	// -------------------
	*commonLogger.LoggerConfig
	// -------------------
	// Vault config
	// -------------------
	// -------------------
	*VaultWrappedConfig
	// Database config
	// -------------------
	*commonPostgres.PostgresConfig
	*commonRedis.RedisConfig
	//*natsCfg.NatsConfig
	// -------------------
	// GRPC service config
	// -------------------
	*GrpcConfig
	// -------------------
	// HD wallet config
	// -------------------
	*MnemonicConfig
	// -------------------
	// Wallet manager config
	// -------------------
	WalletManagerUnloadHotInterval      time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_HOT_INTERVAL" default:"15s"`
	WalletManagerUnloadInterval         time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_INTERVAL" default:"8s"`
	WalletManagerMnemonicPerWalletCount uint8         `envconfig:"WALLET_MANAGER_MNEMONICS_PER_WALLET_COUNT" default:"3"`
}

func (c *Config) GetDefaultHotWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadHotInterval
}

func (c *Config) GetDefaultWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadInterval
}

func (c *Config) GetMnemonicsCountPerWallet() uint8 {
	return c.WalletManagerMnemonicPerWalletCount
}

// Prepare variables to static configuration
func (c *Config) Prepare() error {
	return nil
}

func (c *Config) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
