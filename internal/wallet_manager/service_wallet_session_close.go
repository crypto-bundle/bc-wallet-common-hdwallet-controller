package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"go.uber.org/zap"
)

func (s *Service) CloseWalletSession(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	walletItem, sessionItem, err := s.cacheStoreDataSvc.GetMnemonicWalletSessionInfoByUUID(ctx,
		walletUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	if walletItem != nil && sessionItem != nil {
		session, sessErr := s.closeWalletSession(ctx, walletItem, sessionItem)
		if sessErr != nil {
			return nil, nil, sessErr
		}

		return walletItem, session, nil
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

	sessionItem, err = s.closeWalletSession(ctx, walletItem, sessionItem)
	if err != nil {
		return nil, nil, err
	}

	return walletItem, sessionItem, nil
}

func (s *Service) closeWalletSession(ctx context.Context,
	wallet *entities.MnemonicWallet,
	sessionItem *entities.MnemonicWalletSession,
) (session *entities.MnemonicWalletSession, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedSession, clbErr := s.mnemonicWalletsDataSvc.UpdateWalletSessionStatusBySessionUUID(txStmtCtx,
			wallet.UUID.String(), types.MnemonicWalletSessionStatusClosed)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonic sessions status", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()),
				zap.String(app.MnemonicWalletSessionUUIDTag, sessionItem.UUID))

			return clbErr
		}

		clbErr = s.cacheStoreDataSvc.UnsetWalletSession(txStmtCtx, wallet.UUID.String(), sessionItem.UUID)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()),
				zap.String(app.MnemonicWalletSessionUUIDTag, sessionItem.UUID))

			return clbErr
		}

		session = updatedSession

		return nil
	})
	if err != nil {
		return nil, err
	}

	return
}
