package wallet_manager

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	cfg    configService

	mnemonicWalletsDataSvc mnemonicWalletsDataService
	cacheStoreDataSvc      mnemonicWalletsCacheStoreService
	signReqDataSvc         signRequestDataService

	hdwalletClientSvc hdwallet.HdWalletApiClient

	txStmtManager transactionalStatementManager
}

func NewService(logger *zap.Logger,
	cfg configService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	cacheDataSvc mnemonicWalletsCacheStoreService,
	signReqDataSvc signRequestDataService,
	hdwalletClient hdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) (*Service, error) {
	return &Service{
		logger: logger,
		cfg:    cfg,

		txStmtManager:          txStmtManager,
		hdwalletClientSvc:      hdwalletClient,
		cacheStoreDataSvc:      cacheDataSvc,
		mnemonicWalletsDataSvc: mnemonicWalletDataSrv,
		signReqDataSvc:         signReqDataSvc,
	}, nil
}
