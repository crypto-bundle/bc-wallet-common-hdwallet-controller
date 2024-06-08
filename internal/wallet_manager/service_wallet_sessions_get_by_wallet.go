/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"go.uber.org/zap"
)

func (s *Service) GetWalletSessionsByWalletUUID(ctx context.Context,
	walletUUID string,
) (*entities.MnemonicWallet, []*entities.MnemonicWalletSession, error) {
	walletItem, sessionsList, err := s.cacheStoreDataSvc.GetMnemonicWalletInfoByUUID(ctx, walletUUID)
	if err != nil {
		s.logger.Error("unable get wallet and wallet sessions from cache storage", zap.Error(err))
		// no return - it's ok
		return nil, nil, err
	}

	switch true {
	case walletItem != nil && sessionsList != nil:
		return walletItem, sessionsList, nil

	case walletItem != nil && sessionsList == nil:
		_, sessions, caseErr := s.mnemonicWalletsDataSvc.GetActiveWalletSessionsByWalletUUID(ctx, walletUUID)
		if caseErr != nil {
			return nil, nil, caseErr
		}

		return walletItem, sessions, nil

	case walletItem == nil && sessionsList == nil:
		return s.getWalletAndSessionsFromPersistentStore(ctx, walletUUID)

	case walletItem == nil && sessionsList != nil:
		wallet, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return nil, nil, clbErr
		}

		return wallet, sessionsList, nil
	default:
		return nil, nil, nil
	}
}

func (s *Service) getWalletAndSessionsFromPersistentStore(ctx context.Context,
	walletUUID string,
) (wallet *entities.MnemonicWallet, list []*entities.MnemonicWalletSession, err error) {
	if err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		walletItem, clbErr := s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		wallet = walletItem

		_, sessionsList, clbErr := s.mnemonicWalletsDataSvc.GetActiveWalletSessionsByWalletUUID(ctx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		list = sessionsList

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return
}
