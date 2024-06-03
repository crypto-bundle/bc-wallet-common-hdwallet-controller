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

package controller

import (
	"context"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
)

func (s *ManagerApiClientWrapper) DisableMultipleWallets(ctx context.Context,
	walletUUIDs []string,
) (*DisableWalletsResponse, error) {
	identities := make([]*pbCommon.MnemonicWalletIdentity, len(walletUUIDs))
	for i, _ := range walletUUIDs {
		identities[i] = &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUIDs[i],
		}
	}

	resp, err := s.originGRPCClient.DisableWallets(ctx, &DisableWalletsRequest{
		WalletIdentifiers: identities,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, ErrMissingResponse
	}

	return resp, nil
}
