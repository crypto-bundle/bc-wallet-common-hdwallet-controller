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

package grpc

import (
	"context"
	"fmt"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"github.com/google/uuid"
)

type AccessTokenListForm struct {
	list            []*pbApi.AccessTokenData
	tokensCount     uint
	currentPosition uint
}

func (f *AccessTokenListForm) LoadAndValidate(ctx context.Context,
	list []*pbApi.AccessTokenData,
) (valid bool, err error) {
	if len(list) == 0 {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access tokens list")
	}

	f.list = list
	f.tokensCount = uint(len(list))
	f.currentPosition = 0

	return true, nil
}

func (f *AccessTokenListForm) validateTokenData(
	pbTokenData *pbApi.AccessTokenData,
) (tokenUUID *uuid.UUID, tokenRawData []byte, err error) {
	if pbTokenData == nil {
		return nil, nil, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token data")
	}

	if pbTokenData.AccessTokenIdentifier == nil {
		return nil, nil, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token identity")
	}

	if pbTokenData.AccessTokenData == nil {
		return nil, nil, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token raw data")
	}

	tokenUUIDRaw, err := uuid.Parse(pbTokenData.AccessTokenIdentifier.UUID)
	if err != nil {
		return nil, nil, err
	}

	return &tokenUUIDRaw, pbTokenData.AccessTokenData, nil
}

func (f *AccessTokenListForm) GetCount() uint {
	return f.tokensCount
}

func (f *AccessTokenListForm) Next() (tokenUUID *uuid.UUID, tokenRawData []byte, err error) {
	position := f.currentPosition
	f.currentPosition++

	if position < f.tokensCount {
		return f.validateTokenData(f.list[position])
	}

	return nil, nil, nil
}
