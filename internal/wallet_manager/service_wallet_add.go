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
	requestedAccessTokensCount uint,
) (*entities.MnemonicWallet, []*entities.AccessToken, error) {
	walletUUID := uuid.New()

	resp, err := s.hdWalletClientSvc.GenerateMnemonic(ctx, &hdwallet.GenerateMnemonicRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID.String(),
		},
	})
	if err != nil {
		s.logger.Error("unable to generate new mnemonic", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID.String()))

		return nil, nil, err
	}

	if resp == nil {
		s.logger.Error("missing resp in generate mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID.String()))

		return nil, nil, ErrMissingHdWalletResp
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

	accessTokenList, err := s.generateAccessTokens(ctx, walletUUID, requestedAccessTokensCount)
	if err != nil {
		return nil, nil, err
	}

	return s.saveWalletAndTokens(ctx, toSaveWalletItem, accessTokenList,
		resp.WalletIdentifier, resp.EncryptedMnemonicData)
}

func (s *Service) generateAccessTokens(ctx context.Context,
	walletUUID uuid.UUID,
	requestedAccessTokensCount uint,
) ([]*entities.AccessToken, error) {
	createdAt := time.Now()
	expiredTime := time.Now().
		Add(time.Hour * 24 * 356 * 7).
		Truncate(1 * time.Second)

	tokens := make([]*entities.AccessToken, requestedAccessTokensCount)

	roles := []types.AccessTokenRole{
		types.AccessTokenRoleSigner,
		types.AccessTokenRoleFakeSigner,
		types.AccessTokenRoleReader,
	}
	rolesCount := uint(len(roles)) - 1

	for i := uint(0); i != requestedAccessTokensCount; i++ {
		var role = types.AccessTokenRoleReader
		if i <= rolesCount {
			role = roles[i]
		}

		token, loopErr := s.generateAccessToken(walletUUID, role, createdAt, expiredTime)
		if loopErr != nil {
			return nil, loopErr
		}

		tokens[i] = token
	}

	return tokens, nil
}

func (s *Service) generateAccessToken(walletUUID uuid.UUID,
	tokenRole types.AccessTokenRole,
	createdAt, expiredAt time.Time,
) (*entities.AccessToken, error) {
	tokenUUID := uuid.New()

	tokenStr, err := s.jwtSvc.GenerateJWT(expiredAt, map[string]string{
		"token_uuid":       tokenUUID.String(),
		"token_expired_at": expiredAt.Format(time.DateTime),
		"token_role":       tokenRole.String(),
		"wallet_uuid":      walletUUID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &entities.AccessToken{
		UUID:       tokenUUID,
		Role:       tokenRole,
		WalletUUID: walletUUID,
		RawData:    []byte(tokenStr),
		Hash:       fmt.Sprintf("%x", sha256.Sum256([]byte(tokenStr))),
		CreatedAt:  createdAt,
		ExpiredAt:  expiredAt,
		UpdatedAt:  &createdAt,
	}, nil
}

func (s *Service) saveWalletAndTokens(ctx context.Context,
	walletItem *entities.MnemonicWallet,
	tokenItems []*entities.AccessToken,
	hdWalletInfo *pbCommon.MnemonicWalletIdentity,
	encryptedData []byte,
) (wallet *entities.MnemonicWallet, accessTokens []*entities.AccessToken, err error) {
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

		_, accessTokensList, clbErr := s.accessTokenSvc.AddMultipleAccessTokens(txStmtCtx, tokenItems)
		if clbErr != nil {
			s.logger.Error("unable to save wallet access token items in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletItem.UUID.String()))

			return clbErr
		}

		wallet = savedItem
		accessTokens = accessTokensList

		return nil
	})
	if err != nil {
		s.logger.Error("unable to save new wallet", zap.Error(err))

		return nil, nil, err
	}

	return
}
