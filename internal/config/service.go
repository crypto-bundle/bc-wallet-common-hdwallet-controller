package config

import (
	"context"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

func PrepareVault(ctx context.Context, baseCfgSrv baseConfigService) (*commonVault.Service, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager()
	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{},
		AuthConfig: &commonVaultTokenClient.AuthConfig{},
	}
	err := cfgPreparerSrv.PrepareTo(vaultCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	vaultClientSrv, err := commonVaultTokenClient.NewClient(ctx, vaultCfg)
	if err != nil {
		return nil, err
	}

	vaultSrv, err := commonVault.NewService(ctx, vaultCfg, vaultClientSrv)
	if err != nil {
		return nil, err
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		return nil, err
	}

	return vaultSrv, nil
}

func Prepare(ctx context.Context,
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
	applicationName string,
) (*MangerConfig, *commonVault.Service, error) {
	baseCfgSrv, err := PrepareBaseConfig(ctx, version, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS, applicationName)
	if err != nil {
		return nil, nil, err
	}

	vaultSecretSrv, err := PrepareVault(ctx, baseCfgSrv)
	if err != nil {
		return nil, nil, err
	}

	err = vaultSecretSrv.LoadSecrets(ctx)
	if err != nil {
		return nil, nil, err
	}

	appCfgPreparerSrv := commonConfig.NewConfigManager()
	wrappedConfig := &MangerConfig{}
	err = appCfgPreparerSrv.PrepareTo(wrappedConfig).With(baseCfgSrv, vaultSecretSrv).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	wrappedConfig.BaseConfig = baseCfgSrv

	return wrappedConfig, vaultSecretSrv, nil
}

func PrepareBaseConfig(ctx context.Context,
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
	applicationName string,
) (*commonConfig.BaseConfig, error) {
	flagManagerSrv := commonConfig.NewLdFlagsManager(version, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)

	err := commonConfig.LoadLocalEnvIfDev()
	if err != nil {
		return nil, err
	}

	baseCfgPreparerSrv := commonConfig.NewConfigManager()
	baseCfg := commonConfig.NewBaseConfig(applicationName)
	err = baseCfgPreparerSrv.PrepareTo(baseCfg).With(flagManagerSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	return baseCfg, nil
}
