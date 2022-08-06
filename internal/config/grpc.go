package config

import (
	"fmt"
	"strings"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/app"

	"github.com/kelseyhightower/envconfig"
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
