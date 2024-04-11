package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	"go.uber.org/zap"
)

func (s *Service) EnableWalletByUUID(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, error) {
	var resultItem *entities.MnemonicWallet = nil

	err := s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedItem, clbErr := s.mnemonicWalletsDataSvc.UpdateWalletStatus(txStmtCtx,
			walletUUID, types.MnemonicWalletStatusEnabled)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		clbErr = s.cacheStoreDataSvc.SetMnemonicWalletItem(txStmtCtx, updatedItem)
		if clbErr != nil {
			s.logger.Error("unable to unset mnemonic wallet data from cache store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		resultItem = updatedItem

		return nil
	})
	if err != nil {
		s.logger.Error("unable to enable wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID))

		return nil, err
	}

	return resultItem, nil
}
