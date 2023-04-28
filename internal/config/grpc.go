package config

// GrpcConfig for application
type GrpcConfig struct {
}

// PrepareWith variables with dependencies service-components
func (c *GrpcConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
