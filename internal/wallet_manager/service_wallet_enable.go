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

func (s *Service) EnableWalletByUUID(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, error) {
	var resultItem *entities.MnemonicWallet = nil

	err := s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		updatedItem, clbErr := s.mnemonicWalletsDataSvc.UpdateWalletStatus(txStmtCtx,
			walletUUID, types.MnemonicWalletStatusEnabled)
		if clbErr != nil {
			s.logger.Error("unable to update mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return clbErr
		}

		if updatedItem == nil {
			resultItem = nil

			return nil
		}

		clbErr = s.cacheStoreDataSvc.SetMnemonicWalletItem(txStmtCtx, updatedItem)
		if clbErr != nil {
			s.logger.Error("unable to set mnemonic wallet data in cache store", zap.Error(clbErr),
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
