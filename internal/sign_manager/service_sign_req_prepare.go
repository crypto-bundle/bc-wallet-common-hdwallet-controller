package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

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
			UUID:           uuid.NewString(),
			WalletUUID:     mnemonicUUID,
			SessionUUID:    sessionUUID,
			PurposeUUID:    purposeUUID,
			DerivationPath: []int32{int32(account), int32(change), int32(index)},
			Status:         types.SignRequestStatusCreated,
			CreatedAt:      time.Time{},
			UpdatedAt:      nil,
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

	err = s.eventPublisherSvc.SendSignPreparedEvent(ctx, signReqItem.UUID)
	if err != nil {
		s.logger.Error("unable to broadcast sign request prepared event", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemonicUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, sessionUUID),
			zap.String(app.SignRequestUUIDTag, signReqItem.UUID))

		// no return - it's ok
	}

	return
}

func (s *Service) signPrepare(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
) (signerAddr *pbCommon.DerivationAddressIdentity, err error) {
	resp, err := s.hdWalletClientSvc.LoadDerivationAddress(ctx, &hdwallet.LoadDerivationAddressRequest{
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
			s.logger.Error("unable to load derivation address from hd-wallet",
				zap.Error(ErrUnableDecodeGrpcErrorStatus))

			return nil, err
		}
	}

	return resp.TxOwnerIdentity, nil
}
