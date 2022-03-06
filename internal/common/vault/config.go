package vault

type Config struct {
	Address         string `envconfig:"VAULT_ADDRESS" default:""`
	AppRole         string `envconfig:"VAULT_APP_ROLE" default:""`
	AuthPath        string `envconfig:"VAULT_AUTH_PATH" default:""`
	AuthUsername    string `envconfig:"VAULT_AUTH_USERNAME" default:""`
	AuthPassword    string `envconfig:"VAULT_AUTH_PASSWORD" default:""`
	KubeSATokenPath string `envconfig:"VAULT_KUBE_SA_TOKEN_PATH" default:""`
	DataPath        string `envconfig:"VAULT_DATA_PATH" default:""`
	TransitKey      string `envconfig:"VAULT_TRANSIT_KEY" default:""`
}

// IsEmpty returns true if some prop doesn't initialized except for transit key
func (c *Config) IsEmpty() bool {
	return len(c.AuthPath) == 0 ||
		len(c.AppRole) == 0 ||
		len(c.DataPath) == 0 ||
		len(c.Address) == 0
}
