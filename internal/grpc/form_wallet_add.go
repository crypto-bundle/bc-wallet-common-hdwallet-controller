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

package grpc

import (
	"context"
	"fmt"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"github.com/google/uuid"
)

type WalletAddForm struct {
	TokensCount uint
}

func (f *WalletAddForm) LoadAndValidate(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (valid bool, err error) {
	if req.CreateAccessTokensCount == 0 {
		return false, fmt.Errorf("%s: %d", "Minimal tokens count", 4)
	}

	f.TokensCount = uint(req.CreateAccessTokensCount)
	if req.CreateAccessTokensCount < 4 {
		f.TokensCount = 4
	}

	return true, nil
}

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
		return nil, nil,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token data")
	}

	if pbTokenData.AccessTokenIdentifier == nil {
		return nil, nil,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token identity")
	}

	if pbTokenData.AccessTokenData == nil {
		return nil, nil,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Access token raw data")
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
