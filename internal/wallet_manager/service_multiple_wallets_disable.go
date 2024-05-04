/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DisableWalletsByUUIDList(ctx context.Context,
	walletUUIDs []string,
) (count uint, list []*entities.MnemonicWallet, err error) {
	var updatedItemUUIDs []string

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

		updatedItemUUIDs = make([]string, len(mwIdentities), len(mwIdentities))

		updWalletsCount, updatedWalletItems, clbErr := s.mnemonicWalletsDataSvc.UpdateMultipleWalletsStatusClb(txStmtCtx,
			walletUUIDs, types.MnemonicWalletStatusDisabled, func(idx uint, wallet *entities.MnemonicWallet) error {
				updatedItemUUIDs[idx] = wallet.UUID.String()
				return nil
			})
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
			adapter.MarshallItem)
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

		list = updatedWalletItems
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
			WalletUUID: updatedItemUUIDs[i],
		}

		pbIdentities[i] = pbIdentity
	}

	_, unloadErr := s.hdWalletClientSvc.UnLoadMultipleMnemonics(ctx, &hdwallet.UnLoadMultipleMnemonicsRequest{
		WalletIdentifier: pbIdentities})
	if err != nil {
		s.logger.Error("unable to unload mnemonics", zap.Error(unloadErr),
			zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))

		respStatus, ok := status.FromError(unloadErr)
		if !ok {
			s.logger.Warn("unable to extract response status code", zap.Error(unloadErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))
		}
		switch respStatus.Code() {
		case codes.Internal:
			s.logger.Warn("unable to unload mnemonic wallet", zap.Error(unloadErr),
				zap.Strings(app.MnemonicWalletUUIDTag, updatedItemUUIDs))
		default:
			break
		}
	}

	return
}
