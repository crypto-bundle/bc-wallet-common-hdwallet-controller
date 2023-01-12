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
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"

	natsCfg "github.com/crypto-bundle/bc-wallet-common/pkg/nats/config"
	"github.com/crypto-bundle/bc-wallet-common/pkg/postgres"
	"github.com/crypto-bundle/bc-wallet-common/pkg/redis"

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
	*redis.RedisConfig
	*natsCfg.NatsConfig
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
