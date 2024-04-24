package hdwallet

type HdWalletClientConfig struct {
	ConnectionPath string `envconfig:"HDWALLET_UNIX_SOCKET_PATH" default:"unix:/tmp/hdwallet.sock"`
}

func (c *HdWalletClientConfig) GetConnectionPath() string {
	return c.ConnectionPath
}

func (c *HdWalletClientConfig) Prepare() error {
	return nil
}

func (c *HdWalletClientConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
