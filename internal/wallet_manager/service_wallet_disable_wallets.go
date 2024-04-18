package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
)

func (s *Service) DisableWalletsByUUIDList(ctx context.Context,
	walletUUIDs []string,
) (count uint, list []string, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updWalletsCount, updatedItemUUIDs, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletsStatus(txStmtCtx,
			walletUUIDs, types.MnemonicWalletStatusDisabled)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, walletUUIDs))

			return clbErr
		}

		_, updSessUUIDList, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletSessionStatus(txStmtCtx,
			updatedItemUUIDs, types.MnemonicWalletSessionStatusClosed)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonic sessions status", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))

			return clbErr
		}

		clbErr = s.cacheStoreDataSvc.UnsetMultipleWallets(txStmtCtx, updatedItemUUIDs, updSessUUIDList)
		if clbErr != nil {
			s.logger.Error("unable to unset mnemonics wallet data and sessions from cache store", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs),
				zap.Strings(app.MnemonicWalletSessionUUIDTag, updSessUUIDList))

			return clbErr
		}

		pbIdentities := make([]*common.MnemonicWalletIdentity, updWalletsCount)
		for i := uint(0); i != updWalletsCount; i++ {
			pbIdentity := &common.MnemonicWalletIdentity{
				WalletUUID: updatedItemUUIDs[i],
			}

			pbIdentities[i] = pbIdentity
		}

		_, clbErr = s.hdWalletClientSvc.UnLoadMultipleMnemonics(txStmtCtx, &hdwallet.UnLoadMultipleMnemonicsRequest{
			MnemonicIdentity: pbIdentities})
		if clbErr != nil {
			s.logger.Error("unable to unload mnemonics", zap.Error(clbErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))

			respStatus, ok := status.FromError(clbErr)
			if !ok {
				s.logger.Warn("unable to extract response status code",
					zap.Error(clbErr),
					zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))
			}
			switch respStatus.Code() {
			case codes.Internal:
				s.logger.Warn("unable to unload mnemonic wallet",
					zap.Error(clbErr),
					zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))
			default:

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

	return
}
