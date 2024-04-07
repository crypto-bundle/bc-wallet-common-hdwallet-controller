package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
)

func (s *Service) CheckSession(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
) (isActive bool, err error) {
	wallet, walletSession, err := s.getWalletAndSession(ctx, mnemonicUUID, sessionUUID)
	if err != nil {
		return false, err
	}

	if wallet == nil {
		return false, nil
	}

	if walletSession == nil {
		return false, nil
	}

	return walletSession.IsSessionActive(), nil
}

func (s *Service) getWalletAndSession(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
) (wallet *entities.MnemonicWallet, session *entities.MnemonicWalletSession, err error) {
	wallet, session, err = s.cacheStoreDataSvc.GetMnemonicWalletSessionInfoByUUID(ctx,
		mnemonicUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	if wallet == nil {
		return nil, nil, nil
	}

	if session != nil {
		return wallet, session, nil
	}

	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		wItem, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(txStmtCtx, mnemonicUUID)
		if clbErr != nil {
			return clbErr
		}

		if wItem == nil {
			return nil
		}

		sItem, clbErr := s.mnemonicWalletsDataSvc.GetWalletSessionByUUID(ctx, sessionUUID)
		if clbErr != nil {
			return clbErr
		}

		if sItem == nil {
			return nil
		}

		session = sItem
		wallet = wItem

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return
}
