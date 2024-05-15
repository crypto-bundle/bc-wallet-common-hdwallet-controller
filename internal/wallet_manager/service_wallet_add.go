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
	"crypto/sha256"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) AddNewWallet(ctx context.Context,
	tokensIterator types.AccessTokenListIterator,
) (*entities.MnemonicWallet, error) {
	walletUUID := uuid.New()

	resp, err := s.hdWalletClientSvc.GenerateMnemonic(ctx, &hdwallet.GenerateMnemonicRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID.String(),
		},
	})
	if err != nil {
		s.logger.Error("unable to generate new mnemonic", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID.String()))

		return nil, err
	}

	if resp == nil {
		s.logger.Error("missing resp in generate mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID.String()))

		return nil, ErrMissingHdWalletResp
	}

	saveTime := time.Now()
	toSaveWalletItem := &entities.MnemonicWallet{
		UUID:               walletUUID,
		MnemonicHash:       "",
		Status:             types.MnemonicWalletStatusCreated,
		UnloadInterval:     s.cfg.GetDefaultWalletUnloadInterval(),
		VaultEncrypted:     nil,
		VaultEncryptedHash: "",
		CreatedAt:          saveTime,
		UpdatedAt:          &saveTime,
	}

	accessTokenList, err := s.tokenDataAdapterSvc.Adopt(walletUUID, tokensIterator)
	if err != nil {
		return nil, err
	}

	return s.saveWalletAndToken(ctx, toSaveWalletItem, accessTokenList,
		resp.WalletIdentifier, resp.EncryptedMnemonicData)
}

func (s *Service) saveWalletAndToken(ctx context.Context,
	walletItem *entities.MnemonicWallet,
	tokenItems []*entities.AccessToken,
	hdWalletInfo *pbCommon.MnemonicWalletIdentity,
	encryptedData []byte,
) (wallet *entities.MnemonicWallet, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		walletItem.MnemonicHash = hdWalletInfo.WalletHash
		walletItem.VaultEncryptedHash = fmt.Sprintf("%x", sha256.Sum256(encryptedData))
		walletItem.VaultEncrypted = encryptedData

		savedItem, clbErr := s.mnemonicWalletsDataSvc.AddNewMnemonicWallet(txStmtCtx,
			walletItem)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletItem.UUID.String()))

			return clbErr
		}

		_, _, clbErr = s.accessTokenSvc.AddMultipleAccessTokens(txStmtCtx, tokenItems)
		if clbErr != nil {
			s.logger.Error("unable to save wallet access token items in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletItem.UUID.String()))

			return clbErr
		}

		wallet = savedItem

		return nil
	})
	if err != nil {
		s.logger.Error("unable to save new wallet", zap.Error(err))

		return nil, err
	}

	return
}
