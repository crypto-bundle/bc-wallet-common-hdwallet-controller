/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
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
 *
 */

package wallet_manager

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	cfg    configService

	transitEncryptorSvc encryptService
	appEncryptorSvc     encryptService

	accessTokenSvc         accessTokenDataService
	mnemonicWalletsDataSvc mnemonicWalletsDataService
	cacheStoreDataSvc      mnemonicWalletsCacheStoreService
	signReqDataSvc         signRequestDataService
	jwtSvc                 jwtService
	tokenDataAdapterSvc    tokenDataAdapter

	hdWalletClientSvc hdwallet.HdWalletApiClient

	eventPublisher eventPublisherService

	txStmtManager transactionalStatementManager
}

func NewService(logger *zap.Logger,
	cfg configService,
	transitEncryptorSvc encryptService,
	appEncryptorSvc encryptService,
	accessTokenSvc accessTokenDataService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	cacheDataSvc mnemonicWalletsCacheStoreService,
	signReqDataSvc signRequestDataService,
	jwtSvc jwtService,
	hdWalletClient hdwallet.HdWalletApiClient,
	eventPublisher eventPublisherService,
	txStmtManager transactionalStatementManager,
) *Service {
	return &Service{
		logger: logger,
		cfg:    cfg,

		transitEncryptorSvc: transitEncryptorSvc,
		appEncryptorSvc:     appEncryptorSvc,

		jwtSvc:              jwtSvc,
		tokenDataAdapterSvc: NewTokenDataAdapter(jwtSvc),

		txStmtManager:          txStmtManager,
		hdWalletClientSvc:      hdWalletClient,
		cacheStoreDataSvc:      cacheDataSvc,
		mnemonicWalletsDataSvc: mnemonicWalletDataSrv,
		signReqDataSvc:         signReqDataSvc,
		accessTokenSvc:         accessTokenSvc,

		eventPublisher: eventPublisher,
	}
}
