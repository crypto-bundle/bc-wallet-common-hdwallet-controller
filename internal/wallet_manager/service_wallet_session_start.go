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

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbHdwallet "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func (s *Service) StartSessionForWallet(ctx context.Context,
	wallet *entities.MnemonicWallet,
	accessTokenUUID uuid.UUID,
) (*entities.MnemonicWalletSession, error) {
	session, err := s.startWalletSession(ctx, wallet, accessTokenUUID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) StartWalletSession(ctx context.Context,
	walletUUID string,
	accessTokenUUID uuid.UUID,
) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error) {
	walletItem, err := s.cacheStoreDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		return nil, nil, err
	}

	if walletItem != nil {
		session, sessErr := s.startWalletSession(ctx, walletItem, accessTokenUUID)
		if sessErr != nil {
			return nil, nil, sessErr
		}

		return walletItem, session, nil
	}

	walletItem, err = s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		return nil, nil, err
	}

	sessionItem, err := s.startWalletSession(ctx, walletItem, accessTokenUUID)
	if err != nil {
		return nil, nil, err
	}

	return walletItem, sessionItem, nil
}

func (s *Service) startWalletSession(ctx context.Context,
	wallet *entities.MnemonicWallet,
	accessTokenUUID uuid.UUID,
) (*entities.MnemonicWalletSession, error) {
	var session *entities.MnemonicWalletSession = nil
	err := s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		currentTime := time.Now()
		startedAt := currentTime.Add(s.cfg.GetDefaultWalletSessionDelay())
		expiredAt := startedAt.Add(wallet.UnloadInterval)

		nextSerialNumber, clbErr := s.mnemonicWalletsDataSvc.GetNextWalletSessionNumberByAccessTokenUUID(txStmtCtx,
			accessTokenUUID.String())
		if clbErr != nil {
			return clbErr
		}

		sessionUUID := uuid.New()

		sessionToSave := &entities.MnemonicWalletSession{
			UUID:               sessionUUID.String(),
			MnemonicWalletUUID: wallet.UUID.String(),
			Status:             types.MnemonicWalletSessionStatusPrepared,
			StartedAt:          startedAt,
			ExpiredAt:          expiredAt,

			CreatedAt: currentTime,
			UpdatedAt: nil,
		}

		accessTokenForSession := &entities.AccessTokenWalletSession{
			SerialNumber:   nextSerialNumber,
			AccessTokeUUID: accessTokenUUID,
			SessionUUID:    sessionUUID,
			CreatedAt:      currentTime,
		}

		sessionToSave, clbErr = s.mnemonicWalletsDataSvc.AddNewWalletSession(txStmtCtx,
			sessionToSave)
		if clbErr != nil {
			s.logger.Error("unable to add new wallet session", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()))

			return clbErr
		}

		_, clbErr = s.mnemonicWalletsDataSvc.AddNewWalletSessionAccessTokenItem(txStmtCtx,
			accessTokenForSession)
		if clbErr != nil {
			s.logger.Error("unable to add new wallet session token", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()))

			return clbErr
		}

		session = sessionToSave

		return nil
	})
	if err != nil {
		return nil, err
	}

	timeToLive := s.cfg.GetDefaultWalletSessionDelay() + wallet.UnloadInterval

	_, err = s.hdWalletClientSvc.LoadMnemonic(ctx, &pbHdwallet.LoadMnemonicRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: wallet.UUID.String(),
			WalletHash: wallet.MnemonicHash,
		},
		TimeToLive:            uint64(timeToLive),
		EncryptedMnemonicData: wallet.VaultEncrypted,
	})
	if err != nil {
		s.logger.Error("unable to load mnemonic wallet by hd-wallet service", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()))

		// no return - it's ok
	}

	err = s.cacheStoreDataSvc.SetMnemonicWalletSessionItem(ctx, session)
	if err != nil {
		s.logger.Error("unable to update mnemonics wallets status in persistent store", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()))

		// no return - it's ok
	}

	err = s.eventPublisher.SendSessionStartEvent(ctx, wallet.UUID.String(), session.UUID)
	if err != nil {
		s.logger.Error("unable to broadcast session start event", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()),
			zap.String(app.MnemonicWalletSessionUUIDTag, session.UUID))

		// no return - it's ok
	}

	return session, nil
}
