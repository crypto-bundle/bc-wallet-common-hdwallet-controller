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
	"go.uber.org/zap"
)

func (s *Service) CloseWalletSession(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	walletItem, sessionItem, err := s.getWalletAndSession(ctx, walletUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	if walletItem == nil {
		return nil, nil, nil
	}

	if sessionItem == nil {
		return walletItem, nil, nil
	}

	sessionItem, err = s.closeWalletSession(ctx, walletItem, sessionItem)
	if err != nil {
		return nil, nil, err
	}

	err = s.eventPublisher.SendSessionClosedEvent(ctx, walletItem.UUID.String(), sessionItem.UUID)
	if err != nil {
		s.logger.Error("unable to broadcast session close event", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletItem.UUID.String()),
			zap.String(app.MnemonicWalletSessionUUIDTag, sessionItem.UUID))

		// no return - it's ok
	}

	return walletItem, sessionItem, nil
}

func (s *Service) closeWalletSession(ctx context.Context,
	wallet *entities.MnemonicWallet,
	sessionItem *entities.MnemonicWalletSession,
) (session *entities.MnemonicWalletSession, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedSession, clbErr := s.mnemonicWalletsDataSvc.UpdateWalletSessionStatusBySessionUUID(txStmtCtx,
			sessionItem.UUID, types.MnemonicWalletSessionStatusClosed)
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
