package hdwallet_api

import "fmt"

type HdWalletTronGRPCClientConfig struct {
	TronHdWalletHost           string `envconfig:"BC_WALLET_TRON_HDWALLET_API_SERVICE_HOST" default:"bc-wallet-tron-hdwallet-api"`
	TronHdWalletPort           uint   `envconfig:"BC_WALLET_TRON_HDWALLET_API_SERVICE_PORT" default:"8100"`
	TronHdWalletClientBalancer bool   `envconfig:"BC_WALLET_TRON_HDWALLET_CLIENT_BALANCER" default:"false"`

	TronHdWalletServerAddress string
}

func (o *HdWalletTronGRPCClientConfig) GetHdWalletApiHost() string {
	return o.TronHdWalletHost
}

func (o *HdWalletTronGRPCClientConfig) GetHdWalletApiPort() uint {
	return o.TronHdWalletPort
}

func (o *HdWalletTronGRPCClientConfig) GetHdWalletServerAddress() string {
	return o.TronHdWalletServerAddress
}

func (o *HdWalletTronGRPCClientConfig) IsClientBalancer() bool {
	return o.TronHdWalletClientBalancer
}

func (o *HdWalletTronGRPCClientConfig) Prepare() error {
	if o.TronHdWalletClientBalancer {
		o.TronHdWalletServerAddress = fmt.Sprintf("dns:///%s:%d", o.TronHdWalletHost, o.TronHdWalletPort)
	} else {
		o.TronHdWalletServerAddress = fmt.Sprintf("%s:%d", o.TronHdWalletHost, o.TronHdWalletPort)
	}

	return nil
}

func (c *HdWalletTronGRPCClientConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
