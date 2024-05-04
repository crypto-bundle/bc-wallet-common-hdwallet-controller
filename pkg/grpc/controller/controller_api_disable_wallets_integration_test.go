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
	"testing"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"

	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHdWalletControllerApiClient_DisableWallets(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	createdWallets := make([]*pbCommon.MnemonicWalletIdentity, 3)
	createdWalletsMap := make(map[string]*pbCommon.MnemonicWalletIdentity, len(createdWallets))

	for i := 0; i != 3; i++ {
		createWalletResp, loopErr := client.AddNewWallet(ctx, &AddNewWalletRequest{})
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if createWalletResp == nil {
			t.Fatal("add new wallet empty response")
		}

		walletInfoResp, loopErr := client.GetWalletInfo(ctx, &GetWalletInfoRequest{
			WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
				WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
			},
		})
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if walletInfoResp == nil || walletInfoResp.WalletIdentifier == nil {
			t.Fatal("missing wallet info identifier")
		}

		if walletInfoResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
			t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet status",
				walletInfoResp.WalletStatus, pbCommon.WalletStatus_WALLET_STATUS_CREATED)
		}

		createdWallets[i] = walletInfoResp.WalletIdentifier
		createdWalletsMap[walletInfoResp.WalletIdentifier.WalletUUID] = walletInfoResp.WalletIdentifier
	}

	walletEnableResp, err := client.DisableWallets(ctx, &DisableWalletsRequest{
		WalletIdentifiers: createdWallets,
	})
	if err != nil {
		t.Fatal(err)
	}

	if walletEnableResp == nil {
		t.Fatal("missing enable wallet response")
	}

	for _, walletIdentifier := range createdWallets {
		walletInfoResp, loopErr := client.GetWalletInfo(ctx, &GetWalletInfoRequest{
			WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
				WalletUUID: walletIdentifier.WalletUUID,
			},
		})
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if walletInfoResp == nil || walletInfoResp.WalletIdentifier == nil {
			t.Fatal("missing wallet info identifier")
		}

		if walletInfoResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_DISABLED {
			t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet status",
				walletInfoResp.WalletStatus, pbCommon.WalletStatus_WALLET_STATUS_DISABLED)
		}
	}
}
