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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

func (m *grpcMarshaller) MarshallGetEnabledWallets(
	walletsItems []*entities.MnemonicWallet,
) *pbApi.GetEnabledWalletsResponse {
	walletCount := uint32(len(walletsItems))

	response := &pbApi.GetEnabledWalletsResponse{
		WalletsCount: walletCount,
		WalletsData:  make([]*pbCommon.MnemonicWalletData, walletCount),
		Bookmarks:    make(map[string]uint32, walletCount),
	}

	for i := uint32(0); i != walletCount; i++ {
		item := walletsItems[i]
		if item == nil {
			continue
		}

		walletUUID := item.UUID.String()

		response.WalletsData[i] = &pbCommon.MnemonicWalletData{
			WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
				WalletUUID: walletUUID,
				WalletHash: item.MnemonicHash,
			},
			WalletStatus: pbCommon.WalletStatus(item.Status),
		}
		response.Bookmarks[walletUUID] = i
	}

	return response
}
