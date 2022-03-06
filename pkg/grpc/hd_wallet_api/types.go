package bc_adapter_api

import "go.uber.org/zap"

type loggerService interface {
	NewLoggerEntry(named string) (*zap.Logger, error)
}