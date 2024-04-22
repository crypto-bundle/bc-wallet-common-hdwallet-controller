package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DisableWalletsByUUIDList(ctx context.Context,
	walletUUIDs []string,
) (count uint, list []string, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		mwIdentities, _, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletsByUUIDListAndStatus(txStmtCtx,
			walletUUIDs, []types.MnemonicWalletStatus{
				types.MnemonicWalletStatusEnabled,
				types.MnemonicWalletStatusCreated,
			})
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return clbErr
		}

		if len(mwIdentities) == 0 {
			return nil
		}

		updWalletsCount, updatedItemUUIDs, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletsStatus(txStmtCtx,
			walletUUIDs, types.MnemonicWalletStatusDisabled)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return clbErr
		}

		adapter := newSessionsByWalletDataMapper()
		updatedSessionsCount, _, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletSessionStatusClb(txStmtCtx,
			updatedItemUUIDs, types.MnemonicWalletSessionStatusClosed, []types.MnemonicWalletSessionStatus{
				types.MnemonicWalletSessionStatusPrepared,
			},
			adapter.Marshall)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonic sessions status", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))

			return clbErr
		}

		clbErr = s.cacheStoreDataSvc.UnsetMultipleWallets(txStmtCtx, updatedItemUUIDs)
		if clbErr != nil {
			s.logger.Error("unable to unset mnemonics wallet data from cache store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs),
				zap.Strings(app.MnemonicWalletSessionUUIDTag, adapter.GetSessionsUUIDs()))

			return clbErr
		}

		if updatedSessionsCount > 0 {
			clbErr = s.cacheStoreDataSvc.UnsetMultipleSessions(txStmtCtx, adapter.GetGroupedSessions())
			if clbErr != nil {
				s.logger.Error("unable to unset mnemonics sessions from cache store", zap.Error(clbErr),
					zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs),
					zap.Strings(app.MnemonicWalletSessionUUIDTag, adapter.GetSessionsUUIDs()))

				return clbErr
			}
		}

		list = updatedItemUUIDs
		count = updWalletsCount

		return nil
	})
	if err != nil {
		s.logger.Error("unable to disable wallets", zap.Error(err),
			zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

		return 0, nil, err
	}

	pbIdentities := make([]*common.MnemonicWalletIdentity, count)
	for i := uint(0); i != count; i++ {
		pbIdentity := &common.MnemonicWalletIdentity{
			WalletUUID: list[i],
		}

		pbIdentities[i] = pbIdentity
	}

	_, unloadErr := s.hdWalletClientSvc.UnLoadMultipleMnemonics(ctx, &hdwallet.UnLoadMultipleMnemonicsRequest{
		MnemonicIdentity: pbIdentities})
	if err != nil {
		s.logger.Error("unable to unload mnemonics", zap.Error(unloadErr),
			zap.Strings(app.MnemonicWalletUUIDTag, list))

		respStatus, ok := status.FromError(unloadErr)
		if !ok {
			s.logger.Warn("unable to extract response status code", zap.Error(unloadErr),
				zap.Strings(app.MnemonicWalletUUIDTag, list))
		}
		switch respStatus.Code() {
		case codes.Internal:
			s.logger.Warn("unable to unload mnemonic wallet", zap.Error(unloadErr),
				zap.Strings(app.MnemonicWalletUUIDTag, list))
		default:
			break
		}
	}

	return
}
