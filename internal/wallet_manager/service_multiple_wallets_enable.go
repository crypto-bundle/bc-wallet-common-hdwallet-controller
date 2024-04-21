package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	"go.uber.org/zap"
)

func (s *Service) EnableWalletsByUUIDList(ctx context.Context,
	walletUUIDs []string,
) (count uint, list []string, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		mwIdentities, mwList, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletsByUUIDListAndStatus(txStmtCtx,
			walletUUIDs, []types.MnemonicWalletStatus{
				types.MnemonicWalletStatusDisabled,
				types.MnemonicWalletStatusCreated,
			})
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return clbErr
		}

		if mwList == nil {
			return nil
		}

		updWalletsCount, updatedItemUUIDs, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletsStatus(txStmtCtx,
			mwIdentities, types.MnemonicWalletStatusEnabled)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return clbErr
		}

		if updWalletsCount != uint(len(mwIdentities)) {
			s.logger.Error("something wrong - updated count less than identities list count",
				zap.Error(ErrUpdatedCountNotEqualExpected),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return ErrUpdatedCountNotEqualExpected
		}

		clbErr = s.cacheStoreDataSvc.SetMultipleMnemonicWallets(txStmtCtx, mwList)
		if clbErr != nil {
			s.logger.Error("unable to set mnemonics wallets data to cache store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))

			return clbErr
		}

		count = updWalletsCount
		list = updatedItemUUIDs

		return nil
	})
	if err != nil {
		s.logger.Error("unable to enable wallets", zap.Error(err),
			zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

		return 0, nil, err
	}

	return
}
