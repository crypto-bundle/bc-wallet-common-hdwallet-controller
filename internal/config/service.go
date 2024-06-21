/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package config

import (
	"context"
	"fmt"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
	"log"
	"os"
)

func PrepareLogger(ctx context.Context,
	baseCfgSrv baseConfigService,
) (*commonLogger.LoggerConfig, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager()
	loggerCfg := &commonLogger.LoggerConfig{}

	err := cfgPreparerSrv.PrepareTo(loggerCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	return loggerCfg, nil
}

func PrepareVault(ctx context.Context,
	baseCfgSrv baseConfigService,
	stdLogger *log.Logger,
) (*commonVault.Service, error) {
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

	vaultSrv, err := commonVault.NewService(stdLogger, vaultCfg, vaultClientSrv)
	if err != nil {
		return nil, err
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		return nil, err
	}

	return vaultSrv, nil
}

func PrepareAppCfg(ctx context.Context,
	wrappedBaseCfgSvc *BaseConfigWrapper,
	stdLogger *log.Logger,
) (*MangerConfig, *commonVault.Service, error) {
	baseCfg, loggerCfg, envProcCfg := wrappedBaseCfgSvc.BaseConfig,
		wrappedBaseCfgSvc.LoggerConfig,
		wrappedBaseCfgSvc.ProcessionEnvironmentConfig

	vaultSecretSvc, err := PrepareVault(ctx, baseCfg, stdLogger)
	if err != nil {
		return nil, nil, err
	}

	err = vaultSecretSvc.LoadSecrets(ctx)
	if err != nil {
		return nil, nil, err
	}

	appCfgPreparerSrv := commonConfig.NewConfigManager()
	wrappedConfig := &MangerConfig{}
	err = appCfgPreparerSrv.PrepareTo(wrappedConfig).With(baseCfg,
		loggerCfg, envProcCfg, vaultSecretSvc).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	wrappedConfig.baseAppCfgSvc = baseCfg
	wrappedConfig.loggerCfgSvc = loggerCfg
	wrappedConfig.processingEnvCfgSvc = envProcCfg

	return wrappedConfig, vaultSecretSvc, nil
}

func PrepareBaseConfig(ctx context.Context,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS string,
) (*BaseConfigWrapper, error) {
	appName := fmt.Sprintf(ApplicationManagerNameTpl, os.Getenv(ProcessingNetworkEnvName))

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)
	if err != nil {
		return nil, err
	}

	err = commonConfig.LoadLocalEnvIfDev()
	if err != nil {
		return nil, err
	}

	baseCfgPreparerSvc := commonConfig.NewConfigManager()
	baseCfg := commonConfig.NewBaseConfig(appName)
	err = baseCfgPreparerSvc.PrepareTo(baseCfg).With(flagManagerSvc).Do(ctx)
	if err != nil {
		return nil, err
	}

	loggerConfig, err := PrepareLogger(ctx, baseCfg)
	if err != nil {
		return nil, err
	}

	procEnvConfig := &ProcessionEnvironmentConfig{}
	procEnvCfgPreparerSvc := commonConfig.NewConfigManager()
	err = procEnvCfgPreparerSvc.PrepareTo(procEnvConfig).With(baseCfg).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &BaseConfigWrapper{
		BaseConfig:                  baseCfg,
		LoggerConfig:                loggerConfig,
		ProcessionEnvironmentConfig: procEnvConfig,
	}, nil
}
