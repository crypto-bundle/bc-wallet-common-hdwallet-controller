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
	"google.golang.org/protobuf/types/known/anypb"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type GetAccountForm struct {
	MnemonicWalletUUID    string `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID

	SessionUUID string `valid:"type(string),uuid,required"`

	AccountParameters *anypb.Any `valid:"required"`
}

func (f *GetAccountForm) LoadAndValidate(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (valid bool, err error) {

	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "MnemonicWallet identity")
	}
	f.MnemonicWalletUUID = req.WalletIdentifier.WalletUUID

	if req.SessionIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Session identity")
	}
	f.SessionUUID = req.SessionIdentifier.SessionUUID

	if req.AccountIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Account identity")
	}
	f.AccountParameters = req.AccountIdentifier.Parameters

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
