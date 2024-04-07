package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
)

func (s *Service) GetWalletSessionsByWalletUUID(ctx context.Context,
	walletUUID string,
) (wallet *entities.MnemonicWallet, list []*entities.MnemonicWalletSession, err error) {
	wallet, list, err = s.cacheStoreDataSvc.GetMnemonicWalletInfoByUUID(ctx, walletUUID)
	if err != nil {
		return nil, nil, err
	}

	if wallet != nil {
		return
	}

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
