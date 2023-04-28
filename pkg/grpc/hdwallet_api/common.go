package hdwallet_api

type clientConfigService interface {
	GetHdWalletApiHost() string
	GetHdWalletApiPort() uint
	GetHdWalletServerAddress() string
}
