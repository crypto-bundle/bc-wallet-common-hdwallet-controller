package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
)

func (s *Service) CloseSignRequest(ctx context.Context,
	signReqUUID string,
) error {
	err := s.signReqDataSvc.UpdateSignRequestItemStatus(ctx, signReqUUID,
		types.SignRequestStatusPrepared)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CloseSignRequestBySession(ctx context.Context,
	sessionUUID string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusBySessionUUID(ctx, sessionUUID,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}

func (s *Service) CloseSignRequestByMultipleWallets(ctx context.Context,
	walletUUIDs []string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusByWalletsUUIDList(ctx, walletUUIDs,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}

func (s *Service) CloseSignRequestByWallet(ctx context.Context,
	walletUUID string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusByWalletUUID(ctx, walletUUID,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}
