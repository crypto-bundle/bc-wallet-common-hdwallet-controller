package manager

type processingEnvConfig interface {
	GetProviderName() string
	GetNetworkName() string
}
