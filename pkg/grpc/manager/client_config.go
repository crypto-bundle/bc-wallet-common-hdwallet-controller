package manager

import "fmt"

type HdWalletGRPCClientConfig struct {
	ApiHost string `envconfig:"BC_WALLET_COMMON_HDWALLET_API_SERVICE_HOST" default:"bc-wallet-tron-hdwallet-api"`
	ApiPort uint   `envconfig:"BC_WALLET_COMMON_HDWALLET_API_SERVICE_PORT" default:"8100"`

	serverAddress string
}

func (o *HdWalletGRPCClientConfig) GetHdWalletApiHost() string {
	return o.ApiHost
}

func (o *HdWalletGRPCClientConfig) GetHdWalletApiPort() uint {
	return o.ApiPort
}

func (o *HdWalletGRPCClientConfig) Prepare() error {
	o.serverAddress = fmt.Sprintf("%s:%d", o.ApiHost, o.ApiPort)

	return nil
}

func (o *HdWalletGRPCClientConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
