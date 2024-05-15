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
	"bytes"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/google/uuid"
	"time"
)

type accessTokenListAdapter struct {
	jwtSvc jwtService
}

func (a *accessTokenListAdapter) Adopt(walletUUID uuid.UUID,
	iterator accessTokenListIterator,
) ([]*entities.AccessToken, error) {
	toSaveAccessTokens := make([]*entities.AccessToken, iterator.GetCount())
	position := 0

	saveTime := time.Now()

	for {
		tokenUUID, rawData, loopErr := iterator.Next()
		if loopErr != nil {
			return nil, loopErr
		}

		if rawData == nil {
			break
		}

		uuidFromJwt, expiredAt, err := a.jwtSvc.ExtractFields(rawData)
		if err != nil {
			return nil, err
		}

		if bytes.Compare(tokenUUID[:], uuidFromJwt[:]) != 0 {
			return nil, ErrAccessTokenUUIDMismatched
		}

		tokenItem := &entities.AccessToken{
			UUID:       *tokenUUID,
			WalletUUID: walletUUID,
			RawData:    rawData,
			CreatedAt:  saveTime,
			ExpiredAt:  *expiredAt,
			UpdatedAt:  &saveTime,
		}

		toSaveAccessTokens[position] = tokenItem
		position++
	}

	return toSaveAccessTokens, nil
}

func NewTokenDataAdapter(jwtSvc jwtService) *accessTokenListAdapter {
	return &accessTokenListAdapter{
		jwtSvc: jwtSvc,
	}
}
