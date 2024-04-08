package sign_manager

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	signReqDataSvc signRequestDataService

	hdwalletClientSvc hdwallet.HdWalletApiClient

	txStmtManager transactionalStatementManager
}

func NewService(logger *zap.Logger,
	signReqDataSvc signRequestDataService,
	hdwalletClient hdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) (*Service, error) {
	return &Service{
		logger: logger,

		txStmtManager:     txStmtManager,
		hdwalletClientSvc: hdwalletClient,
		signReqDataSvc:    signReqDataSvc,
	}, nil
}
