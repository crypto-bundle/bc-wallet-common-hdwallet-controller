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
	"log"
	"os"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	"github.com/kelseyhightower/envconfig"
)

const (
	EnvLocal      = "local"
	EnvDev        = "development"
	EnvStaging    = "staging"
	EnvTesting    = "testing"
	EnvProduction = "production"
)

// BaseConfig is config for application base entity like environment, application run mode and etc
type BaseConfig struct {
	// -------------------
	// Application configs
	// -------------------
	// Application environment.
	// allowed: local, dev, testing, staging, production
	Environment string `envconfig:"APP_ENV" default:"development"`
	// Debug mode
	Debug bool `envconfig:"APP_DEBUG" default:"false"`
	// MinimalLogsLevel is a level for setup minimal logger event notification.
	// Allowed: debug, info, warn, error, dpanic, panic, fatal
	MinimalLogsLevel string `envconfig:"LOGGER_LEVEL" default:"debug"`
	StageName        string `envconfig:"APP_STAGE" default:"dev"`

	// ----------------------------
	// Calculated config parameters
	Hostname string
}

// Prepare variables to static configuration
func (c *BaseConfig) Prepare() error {
	err := envconfig.Process(app.ApplicationName, c)
	if err != nil {
		return err
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	c.Hostname = host

	return err
}

// GetHostName ...
func (c *BaseConfig) GetHostName() string {
	return c.Hostname
}

// IsProd ...
func (c *BaseConfig) IsProd() bool {
	return c.Environment == EnvProduction
}

// IsStage ...
func (c *BaseConfig) IsStage() bool {
	return c.Environment == EnvStaging
}

// IsTest ...
func (c *BaseConfig) IsTest() bool {
	return c.Environment == EnvStaging || c.Environment == EnvTesting
}

// IsDev ...
func (c *BaseConfig) IsDev() bool {
	return c.Environment == EnvLocal || c.Environment == EnvDev
}

// IsDebug ...
func (c *BaseConfig) IsDebug() bool {
	return c.Debug
}

// IsLocal ...
func (c *BaseConfig) IsLocal() bool {
	return c.Environment == EnvLocal
}

// GetMinimalLogLevel ...
func (c *BaseConfig) GetMinimalLogLevel() string {
	return c.MinimalLogsLevel
}

// GetStageName is for getting log stage name environment
func (c *BaseConfig) GetStageName() string {
	return c.StageName
}
