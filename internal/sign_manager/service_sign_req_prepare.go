package sign_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PrepareSignRequest(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
	purposeUUID string,
	account, change, index uint32,
) (addr *pbCommon.DerivationAddressIdentity, signReqItem *entities.SignRequest, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		savedItem, clbErr := s.signReqDataSvc.AddSignRequestItem(txStmtCtx, &entities.SignRequest{
			UUID:        uuid.NewString(),
			WalletUUID:  mnemonicUUID,
			SessionUUID: sessionUUID,
			PurposeUUID: purposeUUID,
			Status:      types.SignRequestStatusCreated,
			CreatedAt:   time.Time{},
			UpdatedAt:   nil,
		})
		if clbErr != nil {
			return clbErr
		}

		signOwner, clbErr := s.signPrepare(ctx, mnemonicUUID, account, change, index)
		if clbErr != nil {
			return clbErr
		}

		clbErr = s.signReqDataSvc.UpdateSignRequestItemStatus(txStmtCtx, savedItem.UUID,
			types.SignRequestStatusPrepared)
		if clbErr != nil {
			return clbErr
		}

		addr = signOwner
		signReqItem = savedItem

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *Service) signPrepare(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
) (signerAddr *pbCommon.DerivationAddressIdentity, err error) {
	resp, err := s.hdwalletClientSvc.LoadDerivationAddress(ctx, &hdwallet.LoadDerivationAddressRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentifier: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
	})
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(ErrUnableDecodeGrpcErrorStatus))
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
			return nil, nil

		default:
			s.logger.Error("unable get block by hash from bc-adapter",
				zap.Error(ErrUnableDecodeGrpcErrorStatus))

			return nil, err
		}
	}

	return resp.TxOwnerIdentity, nil
}
