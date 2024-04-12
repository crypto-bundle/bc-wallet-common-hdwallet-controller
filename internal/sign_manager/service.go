package sign_manager

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	signReqDataSvc signRequestDataService

	hdWalletClientSvc hdwallet.HdWalletApiClient

	eventPublisherSvc eventPublisherService

	txStmtManager transactionalStatementManager
}

func NewService(logger *zap.Logger,
	signReqDataSvc signRequestDataService,
	hdWalletClient hdwallet.HdWalletApiClient,
	eventPublisherSvc eventPublisherService,
	txStmtManager transactionalStatementManager,
) *Service {
	return &Service{
		logger: logger,

		hdWalletClientSvc: hdWalletClient,
		signReqDataSvc:    signReqDataSvc,
		eventPublisherSvc: eventPublisherSvc,

		txStmtManager: txStmtManager,
	}
}
