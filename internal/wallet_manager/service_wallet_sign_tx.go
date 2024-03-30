package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"
)

func (s *Service) SignTransaction(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	sessionUUID string,
	transactionData []byte,
) (signer *entities.MnemonicWallet, signedData []byte, err error) {
	var walletSession *entities.MnemonicWalletSession = nil
	var wallet *entities.MnemonicWallet = nil

	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		walletItem, clbErr := s.mnemonicWalletsDataSrv.GetMnemonicWalletByUUID(txStmtCtx, mnemonicUUID)
		if clbErr != nil {
			return clbErr
		}

		if walletItem == nil {
			return nil
		}

		wallet = walletItem

		sessionItem, clbErr := s.mnemonicWalletsDataSrv.GetWalletSessionByUUID(ctx, sessionUUID)
		if clbErr != nil {
			return clbErr
		}

		if sessionItem == nil {
			return nil
		}

		walletSession = sessionItem

		return nil
	}); err != nil {
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

	signResp, err := s.hdwalletClientSvc.SignTransaction(ctx, &hdwallet.SignTransactionRequest{
		MnemonicWalletIdentifier: &common.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentifier: &common.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
		CreatedTxData: transactionData,
	})
	if err != nil {
		return nil, nil, err
	}

	return wallet, signResp.SignedTxData, nil
}
