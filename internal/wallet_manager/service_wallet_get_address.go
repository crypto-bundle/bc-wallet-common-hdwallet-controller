package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetAddress(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	sessionUUID string,
) (ownerWallet *entities.MnemonicWallet, address *string, err error) {
	var walletSession *entities.MnemonicWalletSession = nil
	var wallet *entities.MnemonicWallet = nil

	wallet, walletSession, err = s.getWalletAndSession(ctx, mnemonicUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	if wallet == nil {
		return nil, nil, nil
	}

	if walletSession == nil {
		return wallet, nil, nil
	}

	if !walletSession.IsSessionActive() {
		return wallet, nil, nil
	}

	resp, err := s.hdwalletClientSvc.GetDerivationAddress(ctx, &hdwallet.DerivationAddressRequest{
		MnemonicWalletIdentifier: &common.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentity: &common.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
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
			s.logger.Error("unable to get derivation address",
				zap.Error(ErrUnableDecodeGrpcErrorStatus),
				zap.String(app.MnemonicWalletUUIDTag, mnemonicUUID))

			return nil, nil, err
		}
	}

	return wallet, &resp.AddressIdentity.Address, nil
}
