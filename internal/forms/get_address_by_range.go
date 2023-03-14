/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

package forms

import (
	"context"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	protoTypes "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/types"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/metadata"
)

type DerivationAddressByRangeForm struct {
	WalletUUID string `valid:"type(string),uuid,required"`

	AccountIndex     uint32 `valid:"type(uint32)"`
	InternalIndex    uint32 `valid:"type(uint32)"`
	AddressIndexFrom uint32 `valid:"type(uint32)"`
	AddressIndexTo   uint32 `valid:"type(uint32)"`
}

func (f *DerivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, ErrUnableReadGrpcMetadata
	}

	walletHeaders := headers.Get(protoTypes.WalletUUIDHeaderName)

	if len(walletHeaders) == 0 {
		return false, ErrUnableGetWalletUUIDFromMetadata
	}

	f.WalletUUID = walletHeaders[0]

	f.AccountIndex = req.AccountIndex
	f.InternalIndex = req.InternalIndex
	f.AddressIndexFrom = req.AddressIndexFrom
	f.AddressIndexTo = req.AddressIndexTo

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
