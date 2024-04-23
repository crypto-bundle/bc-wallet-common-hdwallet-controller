package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ExecuteSignRequest(ctx context.Context,
	signReqItem *entities.SignRequest,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	if signReqItem.DerivationPath == nil {
		return nil, nil, ErrMissingDerivationPathField
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		clbErr := s.signReqDataSvc.UpdateSignRequestItemStatus(txStmtCtx, signReqItem.UUID,
			types.SignRequestStatusSigned)
		if clbErr != nil {
			return clbErr
		}

		accountIdx, internalIdx, addrIdx := uint32(signReqItem.DerivationPath[0]),
			uint32(signReqItem.DerivationPath[1]),
			uint32(signReqItem.DerivationPath[2])

		signerAddr, signedData, clbErr = s.signTransaction(txStmtCtx, signReqItem.WalletUUID,
			accountIdx, internalIdx, addrIdx, transactionData)
		if clbErr != nil {
			return clbErr
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *Service) signTransaction(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	signResp, err := s.hdWalletClientSvc.SignData(ctx, &hdwallet.SignDataRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentifier: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
		DataForSign: transactionData,
	})
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(ErrUnableDecodeGrpcErrorStatus))
			return nil, nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
			return nil, nil, nil

		default:
			s.logger.Error("unable get block by hash from bc-adapter",
				zap.Error(ErrUnableDecodeGrpcErrorStatus))

			return nil, nil, err
		}
	}

	return signResp.TxOwnerIdentity, signResp.SignedData, nil
}
