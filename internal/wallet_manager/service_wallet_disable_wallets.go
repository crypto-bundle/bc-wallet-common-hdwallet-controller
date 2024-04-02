package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

	"go.uber.org/zap"
)

func (s *Service) DisableWalletsByUUIDList(ctx context.Context,
	walletUUIDs []string,
) (uint, []string, error) {
	var resultItem *entities.MnemonicWallet = nil

	item, err := s.mnemonicWalletsDataSrv.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		s.logger.Error("unable to get wallet by uuid", zap.Error(err))

		return 0, nil, err
	}

	if item == nil {
		return 0, nil, nil
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedCount, updatedItem, clbErr := s.mnemonicWalletsDataSrv.UpdateMultipleWalletsStatus(txStmtCtx,
			walletUUIDs, types.MnemonicWalletStatusDisabled)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		clbErr = s.mnemonicWalletsDataSrv.UpdateWalletSessionStatusByWalletUUID(txStmtCtx,
			walletUUID, types.MnemonicWalletSessionStatusClosed)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonic sessions status", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		clbErr = s.cacheStoreDataSvc.FullUnsetMnemonicWallet(txStmtCtx, walletUUID)
		if clbErr != nil {
			s.logger.Error("unable to unset mnemonic wallet data from cache store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		resultItem = updatedItem

		_, clbErr = s.hdwalletClientSvc.UnLoadMnemonic(txStmtCtx, &hdwallet.UnLoadMnemonicRequest{
			MnemonicIdentity: &common.MnemonicWalletIdentity{
				WalletUUID: walletUUID,
			}})
		if clbErr != nil {
			s.logger.Error("unable to unload mnemonics", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))
		}

		return nil
	})
	if err != nil {
		s.logger.Error("unable to disable wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID))

		return nil, err
	}

	return resultItem, nil
}
