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
	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

type DerivationAddressByRangeForm struct {
	WalletUUID            string `valid:"type(string),uuid,required"`
	WalletUUIDRaw         uuid.UUID
	MnemonicWalletUUID    string `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID

	AccountIndex     uint32 `valid:"type(uint32),int,required"`
	InternalIndex    uint32 `valid:"type(uint32),int,required"`
	AddressIndexFrom uint32 `valid:"type(uint32),int,required"`
	AddressIndexTo   uint32 `valid:"type(uint32),int,required"`
}

func (f *DerivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	f.WalletUUID = req.WalletIdentity.WalletUUID
	f.MnemonicWalletUUID = req.MnemonicIdentity.WalletUUID

	f.AccountIndex = req.AccountIndex
	f.InternalIndex = req.InternalIndex
	f.AddressIndexFrom = req.AddressIndexFrom
	f.AddressIndexTo = req.AddressIndexTo

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	walletUUIDRaw, err := uuid.Parse(f.WalletUUID)
	if err != nil {
		return false, err
	}
	f.WalletUUIDRaw = walletUUIDRaw

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
