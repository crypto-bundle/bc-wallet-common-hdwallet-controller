package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
)

type signRequestDataService interface {
	AddSignRequestItem(ctx context.Context,
		toSaveItem *entities.SignRequest,
	) (*entities.SignRequest, error)
	UpdateSignRequestItemStatus(ctx context.Context,
		signReqUUID string,
		newStatus types.SignRequestStatus,
	) error
	UpdateSignRequestItemStatusBySessionUUID(ctx context.Context,
		sessionUUID string,
		newStatus types.SignRequestStatus,
	) (uint, []*entities.SignRequest, error)
	UpdateSignRequestItemStatusByWalletsUUIDList(ctx context.Context,
		walletUUIDs []string,
		newStatus types.SignRequestStatus,
	) (uint, []*entities.SignRequest, error)
	UpdateSignRequestItemStatusByWalletUUID(ctx context.Context,
		walletUUID string,
		newStatus types.SignRequestStatus,
	) (uint, []*entities.SignRequest, error)
	GetSignRequestItemByUUIDAndStatus(ctx context.Context,
		signReqUUID string,
		status types.SignRequestStatus,
	) (*entities.SignRequest, error)
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}
