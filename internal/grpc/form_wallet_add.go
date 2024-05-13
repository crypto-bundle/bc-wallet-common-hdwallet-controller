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

type AddWalletForm struct {
	TokenUUID uuid.UUID
	TokenData []byte
}

func (f *AddWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (valid bool, err error) {
	if req.TokenData == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token data")
	}

	if req.AccessTokenIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token identity")
	}

	tokenUUIDRaw, err := uuid.Parse(req.AccessTokenIdentifier.UUID)
	if err != nil {
		return false, err
	}

	f.TokenUUID = tokenUUIDRaw
	f.TokenData = req.TokenData

	return true, nil
}
