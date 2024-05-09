package hdwallet

type HdWalletClientConfig struct {
	ConnectionPath             string `envconfig:"HDWALLET_UNIX_SOCKET_DIR_PATH" default:"unix:/tmp"`
	UnitSocketFileNameTemplate string `envconfig:"HDWALLET_UNIX_SOCKET_FILE_TEMPLATE" default:"hdwallet_"`
}

func (c *HdWalletClientConfig) GetConnectionPath() string {
	return c.ConnectionPath
}

func (c *HdWalletClientConfig) GetUnixFileNameTemplate() string {
	return c.UnitSocketFileNameTemplate
}

func (c *HdWalletClientConfig) Prepare() error {
	return nil
}

func (c *HdWalletClientConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
