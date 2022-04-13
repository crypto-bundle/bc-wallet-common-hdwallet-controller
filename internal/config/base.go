package config

import (
	"log"
	"os"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/app"

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
