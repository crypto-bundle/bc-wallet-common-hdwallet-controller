/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package config

import (
	"context"
	"log"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"

	"github.com/joho/godotenv"
)

func Prepare(ctx context.Context,
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
) (*Config, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager()
	vaultCfg := &commonVault.Config{}
	err := cfgPreparerSrv.PrepareTo(vaultCfg).Do(ctx)
	if err != nil {
		return nil, err
	}

	flagManagerSrv := commonConfig.NewLdFlagsManager(version, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)

	baseCfg := commonConfig.NewBaseConfig()
	err = cfgPreparerSrv.PrepareTo(baseCfg).With(flagManagerSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	if baseCfg.IsDev() {
		loadErr := godotenv.Load(".env")
		if loadErr != nil {
			log.Fatal(loadErr)
		}
	}

	wrappedConfig := &Config{}
	err = cfgPreparerSrv.With(baseCfg).PrepareTo(wrappedConfig).Do(ctx)
	if err != nil {
		return nil, err
	}

	wrappedConfig.BaseConfig = baseCfg

	return wrappedConfig, nil
}
