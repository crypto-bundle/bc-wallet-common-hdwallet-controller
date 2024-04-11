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

	mnemonicWalletsDataSvc mnemonicWalletsDataService
	cacheStoreDataSvc      mnemonicWalletsCacheStoreService
	signReqDataSvc         signRequestDataService

	hdwalletClientSvc hdwallet.HdWalletApiClient

	txStmtManager transactionalStatementManager
}

func NewService(logger *zap.Logger,
	cfg configService,
	transitEncryptorSvc encryptService,
	appEncryptorSvc encryptService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	cacheDataSvc mnemonicWalletsCacheStoreService,
	signReqDataSvc signRequestDataService,
	hdwalletClient hdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) *Service {
	return &Service{
		logger: logger,
		cfg:    cfg,

		transitEncryptorSvc: transitEncryptorSvc,
		appEncryptorSvc:     appEncryptorSvc,

		txStmtManager:          txStmtManager,
		hdwalletClientSvc:      hdwalletClient,
		cacheStoreDataSvc:      cacheDataSvc,
		mnemonicWalletsDataSvc: mnemonicWalletDataSrv,
		signReqDataSvc:         signReqDataSvc,
	}
}
