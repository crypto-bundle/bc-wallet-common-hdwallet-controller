package config

import (
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

// GrpcConfig for application
type GrpcConfig struct {
	// BindRaw port string, default "8080"
	BindRaw string `envconfig:"API_GRPC_PORT"`

	// ----------------------------
	// Calculated config parameters
	Bind string
}

func (c *GrpcConfig) GetBindPort() string {
	return c.Bind
}

// Prepare variables to static configuration
func (c *GrpcConfig) Prepare() error {
	err := envconfig.Process(app.APIConfigPrefix, c)
	if err != nil {
		return err
	}

	c.Bind = fmt.Sprintf(":%s", strings.TrimLeft(c.BindRaw, ":"))

	return err
}
