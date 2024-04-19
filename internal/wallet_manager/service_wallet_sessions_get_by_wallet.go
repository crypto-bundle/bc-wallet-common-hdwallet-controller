package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
)

func (s *Service) GetWalletSessionsByWalletUUID(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, []*entities.MnemonicWalletSession, error) {
	walletItem, sessionsList, err := s.cacheStoreDataSvc.GetMnemonicWalletInfoByUUID(ctx, walletUUID)
	if err != nil {
		return nil, nil, err
	}

	switch true {
	case walletItem != nil && sessionsList != nil:
		return walletItem, sessionsList, nil

	case walletItem != nil && sessionsList == nil:
		_, sessions, caseErr := s.mnemonicWalletsDataSvc.GetActiveWalletSessionsByWalletUUID(ctx, walletUUID)
		if caseErr != nil {
			return nil, nil, caseErr
		}

		return walletItem, sessions, nil

	case walletItem == nil && sessionsList == nil:
		return s.getWalletAndSessionsFromPersistentStore(ctx, walletUUID)

	case walletItem == nil && sessionsList != nil:
		wallet, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return nil, nil, clbErr
		}

		return wallet, sessionsList, nil
	default:
		return nil, nil, nil
	}
}

func (s *Service) getWalletAndSessionsFromPersistentStore(ctx context.Context,
	walletUUID string,
) (wallet *entities.MnemonicWallet, list []*entities.MnemonicWalletSession, err error) {
	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		walletItem, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		wallet = walletItem

		_, sessionsList, clbErr := s.mnemonicWalletsDataSvc.GetActiveWalletSessionsByWalletUUID(ctx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		list = sessionsList

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return
}
