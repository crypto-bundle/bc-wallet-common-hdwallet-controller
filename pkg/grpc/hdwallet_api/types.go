package hdwallet_api

import (
	"go.uber.org/zap"
)

type clientConfig interface {
	GetHdWalletApiHost() string
	GetHdWalletApiPort() string
	GetHdWalletServerAddress() string
}

type loggerService interface {
	NewLoggerEntry(named string) (*zap.Logger, error)
}
