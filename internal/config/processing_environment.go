package config

type ProcessionEnvironmentConfig struct {
	ProcessingProvider string `envconfig:"PROCESSING_PROVIDER" default:"cryptobundle"`
	ProcessingNetwork  string `envconfig:"PROCESSING_NETWORK" default:"tron"`

	baseAppCfgSrv baseConfigService
}

// GetProviderName is for getting event filter by processing provider
func (c *ProcessionEnvironmentConfig) GetProviderName() string {
	return c.ProcessingProvider
}

// GetNetworkName is for getting event filter by processing network
func (c *ProcessionEnvironmentConfig) GetNetworkName() string {
	return c.ProcessingNetwork
}
