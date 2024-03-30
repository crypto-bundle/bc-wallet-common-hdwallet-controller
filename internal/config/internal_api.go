package config

type InternalApiConfig struct {
	InternalApiSocketPath string `envconfig:"BC_WALLET_COMMON_HDWALLET_INTERNAL_API_SOCKET_PATH" default:"/tmp/hdwallet_unix.sock"`

	baseAppCfgSrv baseConfigService
}

//nolint:funlen // its ok
func (c *InternalApiConfig) Prepare() error {
	return nil
}

//nolint:funlen // its ok
func (c *InternalApiConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
