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

package access_tokens

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrMissingTokenUUIDIdentity = errors.New("missing uuid in token data")
	ErrMismatchedUUIDIdentity   = errors.New("token identity mismatched")
)

const (
	TokenUUIDLabel    = "token_uuid"
	TokenExpiredLabel = "token_expired_at"
)

type accessTokenManager struct {
	logger  *zap.Logger
	JWTSvc  jwtService
	dataSvc accessTokenDataService
}

func (m *accessTokenManager) ValidateAccessToken(ctx context.Context,
	tokenData []byte,
) (*uuid.UUID, *time.Time, error) {
	return m.validateAccessToken(ctx, tokenData)
}

func (m *accessTokenManager) ValidateSaveAccessToken(ctx context.Context,
	walletUUID uuid.UUID,
	tokenUUID uuid.UUID,
	tokenData []byte,
) (result *entities.AccessToken, err error) {
	accessToken, err := m.extractAccessTokenFromData(ctx, tokenUUID, tokenData)
	if err != nil {
		return nil, err
	}

	accessToken.WalletUUID = walletUUID

	return m.dataSvc.AddNewAccessToken(ctx, accessToken)
}

func (m *accessTokenManager) SaveAccessToken(ctx context.Context,
	token *entities.AccessToken,
) (result *entities.AccessToken, err error) {
	return m.dataSvc.AddNewAccessToken(ctx, token)
}

func (m *accessTokenManager) GetAccessTokenByUUID(ctx context.Context,
	tokenUUID string,
) (result *entities.AccessToken, err error) {
	return m.dataSvc.GetAccessTokenInfoByUUID(ctx, tokenUUID)
}

func (m *accessTokenManager) ExtractAccessTokenFromData(ctx context.Context,
	tokenUUID uuid.UUID,
	tokenData []byte,
) (result *entities.AccessToken, err error) {
	return m.extractAccessTokenFromData(ctx, tokenUUID, tokenData)
}

func (m *accessTokenManager) validateAccessToken(ctx context.Context,
	tokenData []byte,
) (*uuid.UUID, *time.Time, error) {
	data, err := m.JWTSvc.GetTokenData(string(tokenData))
	if err != nil {
		return nil, nil, err
	}

	tokenUUIDStr, isExist := data[TokenUUIDLabel]
	if !isExist {
		return nil, nil, ErrMissingTokenUUIDIdentity
	}

	tokenUUIDRaw, err := uuid.Parse(tokenUUIDStr)
	if err != nil {
		return nil, nil, err
	}

	tokenExpiredAtStr, isExist := data[TokenExpiredLabel]
	if !isExist {
		return nil, nil, ErrMissingTokenUUIDIdentity
	}

	expiredAt, err := time.Parse(time.Layout, tokenExpiredAtStr)
	if err != nil {
		return nil, nil, err
	}

	return &tokenUUIDRaw, &expiredAt, nil
}

func (m *accessTokenManager) extractAccessTokenFromData(ctx context.Context,
	tokenUUID uuid.UUID,
	tokenData []byte,
) (result *entities.AccessToken, err error) {
	tokenUUIDRaw, expiredAt, err := m.validateAccessToken(ctx, tokenData)
	if err != nil {
		return nil, err
	}

	cmp := bytes.Compare(tokenUUID[:], tokenUUIDRaw[:])
	if cmp != 0 {
		return nil, ErrMismatchedUUIDIdentity
	}

	nowTime := time.Now()

	return &entities.AccessToken{
		UUID:       *tokenUUIDRaw,
		WalletUUID: uuid.Nil,
		RawData:    tokenData,
		CreatedAt:  nowTime,
		ExpiredAt:  *expiredAt,
		UpdatedAt:  &nowTime,
	}, nil
}

func NewTokenManager(
	logger *zap.Logger,
	JWTSvc jwtService,
	dataSvc accessTokenDataService,
) *accessTokenManager {
	return &accessTokenManager{
		logger:  logger,
		JWTSvc:  JWTSvc,
		dataSvc: dataSvc,
	}
}
