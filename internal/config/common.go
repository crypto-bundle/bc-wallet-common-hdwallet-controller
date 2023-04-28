package config

import (
	"time"

	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type baseConfigService interface {
	GetHostName() string
	GetEnvironmentName() string
	IsProd() bool
	IsStage() bool
	IsTest() bool
	IsDev() bool
	IsDebug() bool
	IsLocal() bool
	GetStageName() string
	GetApplicationPID() int
	GetApplicationName() string
	SetApplicationName(appName string)
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}
