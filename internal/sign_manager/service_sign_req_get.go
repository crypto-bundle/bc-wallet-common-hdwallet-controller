package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
)

func (s *Service) GetActiveSignRequest(ctx context.Context,
	signReqUUID string,
) (*entities.SignRequest, error) {
	signReqItem, err := s.signReqDataSvc.GetSignRequestItemByUUIDAndStatus(ctx, signReqUUID,
		types.SignRequestStatusPrepared)
	if err != nil {
		return nil, err
	}

	if signReqItem == nil {
		return nil, nil
	}

	return nil, nil
}
