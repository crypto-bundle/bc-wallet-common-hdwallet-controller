package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
)

func (s *Service) GetWalletSessionInfo(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	walletItem, sessionItem, err := s.cacheStoreDataSvc.GetMnemonicWalletSessionInfoByUUID(ctx, walletUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	if walletItem != nil && sessionItem != nil {
		return walletItem, sessionItem, nil
	}

	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		wallet, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		session, clbErr := s.mnemonicWalletsDataSvc.GetWalletSessionByUUID(ctx, sessionUUID)
		if clbErr != nil {
			return clbErr
		}

		walletItem = wallet
		sessionItem = session

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return walletItem, sessionItem, nil
}
