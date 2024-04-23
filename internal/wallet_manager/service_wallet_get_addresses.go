package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetAddressesByRange(ctx context.Context,
	mnemonicUUID string,
	addrRanges []*pbCommon.RangeRequestUnit,
) (list []*pbCommon.DerivationAddressIdentity, err error) {

	resp, err := s.hdWalletClientSvc.GetDerivationAddressByRange(ctx, &hdwallet.DerivationAddressByRangeRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		Ranges: addrRanges,
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
			s.logger.Error("unable to get derivation address",
				zap.Error(ErrUnableDecodeGrpcErrorStatus),
				zap.String(app.MnemonicWalletUUIDTag, mnemonicUUID))

			return nil, err
		}
	}

	return resp.AddressIdentities, nil
}
