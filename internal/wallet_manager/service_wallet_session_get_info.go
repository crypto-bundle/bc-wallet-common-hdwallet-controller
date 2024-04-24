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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
)

func (s *Service) GetWalletSessionInfo(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	return s.getWalletAndSession(ctx, walletUUID, sessionUUID)
}

func (s *Service) getWalletAndSession(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	walletItem, sessionItem, err := s.cacheStoreDataSvc.GetMnemonicWalletSessionInfoByUUID(ctx, walletUUID, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	switch true {
	case walletItem != nil && sessionItem != nil:
		return walletItem, sessionItem, nil

	case walletItem != nil && sessionItem == nil:
		session, clbErr := s.mnemonicWalletsDataSvc.GetWalletSessionByUUID(ctx, sessionUUID)
		if clbErr != nil {
			return nil, nil, clbErr
		}

		return walletItem, session, nil

	case walletItem == nil && sessionItem == nil:
		return s.getWalletAndSessionFromPersistentStore(ctx, walletUUID, sessionUUID)

	case walletItem == nil && sessionItem != nil:
		wallet, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return nil, nil, clbErr
		}

		return wallet, sessionItem, nil
	default:
		return nil, nil, nil
	}
}

func (s *Service) getWalletAndSessionFromPersistentStore(ctx context.Context,
	walletUUID, sessionUUID string,
) (walletItem *entities.MnemonicWallet, sessionItem *entities.MnemonicWalletSession, err error) {
	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		wallet, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(txStmtCtx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		session, clbErr := s.mnemonicWalletsDataSvc.GetWalletSessionByUUID(txStmtCtx, sessionUUID)
		if clbErr != nil {
			return clbErr
		}

		walletItem = wallet
		sessionItem = session

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return
}
