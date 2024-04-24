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
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	"github.com/asaskevich/govalidator"
)

type WalletsIdentitiesForm struct {
	WalletUUIDs []string `valid:"type([]string),required"`
}

func (f *WalletsIdentitiesForm) LoadAndValidate(
	list []*pbCommon.MnemonicWalletIdentity,
) (valid bool, err error) {
	count := len(list)
	if count == 0 {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identities")
	}

	f.WalletUUIDs = make([]string, count)

	for i, v := range list {
		if !govalidator.IsUUID(v.WalletUUID) {
			return false, fmt.Errorf("%s does not validate as %s", v.WalletUUID, "UUID")
		}

		f.WalletUUIDs[i] = v.WalletUUID
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
