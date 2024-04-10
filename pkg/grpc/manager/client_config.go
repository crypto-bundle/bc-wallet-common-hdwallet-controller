package manager

import (
	"fmt"
	"os"
	"strconv"
)

const (
	hostTemplate = "BC_WALLET_%s_HDWALLET_CONTROLLER_SERVICE_HOST"
	portTemplate = "BC_WALLET_%s_HDWALLET_CONTROLLER_SERVICE_PORT"
)

type HdWalletControllerGRPCClientConfig struct {

	// dependencies
	processingEnvCfg processingEnvConfig

	// calculated
	serverAddress string
	apiHost       string
	apiPort       uint
}

func (o *HdWalletControllerGRPCClientConfig) GetHdWalletApiHost() string {
	return o.apiHost
}

func (o *HdWalletControllerGRPCClientConfig) GetHdWalletApiPort() uint {
	return o.apiPort
}

func (o *HdWalletControllerGRPCClientConfig) Prepare() error {
	o.apiHost = os.Getenv(fmt.Sprintf(hostTemplate, o.processingEnvCfg.GetNetworkName()))
	portRaw := os.Getenv(fmt.Sprintf(portTemplate, o.processingEnvCfg.GetNetworkName()))
	port, err := strconv.ParseUint(portRaw, 10, 0)
	if err != nil {
		return err
	}
	o.apiPort = uint(port)

	o.serverAddress = fmt.Sprintf("%s:%d", o.apiHost, o.apiPort)

	return nil
}

func (o *HdWalletControllerGRPCClientConfig) PrepareWith(cfgSrvList ...interface{}) error {
	for _, cfgSrv := range cfgSrvList {
		switch castedCfg := cfgSrv.(type) {
		case processingEnvConfig:
			o.processingEnvCfg = castedCfg
		default:
			continue
		}
	}

	return nil
}
