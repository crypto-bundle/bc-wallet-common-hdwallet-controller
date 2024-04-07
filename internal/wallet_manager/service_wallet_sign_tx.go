package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SignTransactionWithWallet(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
	account, change, index uint32,
	transactionData []byte,
) (*entities.MnemonicWallet, []byte, error) {
	var walletSession *entities.MnemonicWalletSession = nil
	var wallet *entities.MnemonicWallet = nil

	wallet, walletSession, err := s.getWalletAndSession(ctx, mnemonicUUID, sessionUUID)
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

	_, signedData, err := s.signTransaction(ctx, mnemonicUUID, account, change, index, transactionData)
	if err != nil {
		return nil, nil, err
	}

	return wallet, signedData, nil
}

func (s *Service) SignTransaction(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	return s.signTransaction(ctx, mnemonicUUID, account, change, index, transactionData)
}

func (s *Service) signTransaction(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	signResp, err := s.hdwalletClientSvc.SignTransaction(ctx, &hdwallet.SignTransactionRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentifier: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
		CreatedTxData: transactionData,
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

	return signResp.TxOwnerIdentity, signResp.SignedTxData, nil
}
