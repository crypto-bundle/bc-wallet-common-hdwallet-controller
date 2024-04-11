package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
)

func (s *Service) DisableWalletByUUID(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, error) {
	var resultItem *entities.MnemonicWallet = nil

	item, err := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		s.logger.Error("unable to get wallet by uuid", zap.Error(err))

		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedItem, clbErr := s.mnemonicWalletsDataSvc.UpdateWalletStatus(txStmtCtx,
			walletUUID, types.MnemonicWalletStatusDisabled)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		clbErr = s.mnemonicWalletsDataSvc.UpdateWalletSessionStatusByWalletUUID(txStmtCtx,
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
