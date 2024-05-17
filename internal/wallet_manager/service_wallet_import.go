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
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ImportWallet(ctx context.Context,
	importedData []byte,
	requestedAccessTokensCount uint,
) (*entities.MnemonicWallet, error) {
	decryptedData, err := s.transitEncryptorSvc.Decrypt(importedData)
	if err != nil {
		return nil, err
	}

	encryptedMnemonicData, err := s.appEncryptorSvc.Encrypt(decryptedData)
	if err != nil {
		return nil, err
	}

	mnemonicHash := fmt.Sprintf("%x", sha256.Sum256(decryptedData))
	vaultEncryptedHash := fmt.Sprintf("%x", sha256.Sum256(encryptedMnemonicData))
	walletUUID := uuid.New()

	resp, err := s.hdWalletClientSvc.ValidateMnemonic(ctx, &hdwallet.ValidateMnemonicRequest{
		WalletIdentifier: &common.MnemonicWalletIdentity{
			WalletUUID: walletUUID.String(),
		},
		MnemonicData: encryptedMnemonicData,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		s.logger.Error("missing resp in load mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID.String()))

		return nil, ErrMissingHdWalletResp
	}

	if !resp.IsValid {
		return nil, ErrMnemonicIsNotValid
	}

	toSaveItem := &entities.MnemonicWallet{
		UUID:               walletUUID,
		MnemonicHash:       mnemonicHash,
		Status:             types.MnemonicWalletStatusCreated,
		UnloadInterval:     s.cfg.GetDefaultWalletUnloadInterval(),
		VaultEncrypted:     encryptedMnemonicData,
		VaultEncryptedHash: vaultEncryptedHash,
		CreatedAt:          time.Now(),
		UpdatedAt:          nil,
	}

	accessTokenList, err := s.generateAccessTokens(ctx, walletUUID, requestedAccessTokensCount)
	if err != nil {
		return nil, err
	}

	return s.saveWalletAndTokens(ctx, toSaveItem, accessTokenList,
		resp.WalletIdentifier, encryptedMnemonicData)
}
