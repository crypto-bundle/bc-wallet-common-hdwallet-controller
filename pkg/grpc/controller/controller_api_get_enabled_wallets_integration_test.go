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

func TestHdWalletControllerApiClient_GetEnabledWallets(t *testing.T) {
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

	createdWallets := make(map[string]*pbCommon.MnemonicWalletIdentity, 3)

	for i := 0; i != 3; i++ {
		createWalletResp, loopErr := client.AddNewWallet(ctx, &AddNewWalletRequest{})
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if createWalletResp == nil {
			t.Fatal("add new wallet empty response")
		}

		walletEnableResp, loopErr := client.EnableWallet(ctx, &EnableWalletRequest{
			WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
				WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		if walletEnableResp == nil {
			t.Fatal("missing enable wallet response")
		}

		createdWallets[createWalletResp.WalletIdentifier.WalletUUID] = createWalletResp.WalletIdentifier
	}

	walletsListResp, err := client.GetEnabledWallets(ctx, &GetEnabledWalletsRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if walletsListResp == nil {
		t.Fatal("missing wallet list response")
	}

	if walletsListResp.WalletIdentities == nil {
		t.Fatal("missing wallet identities list in response")
	}

	if len(walletsListResp.WalletIdentities) == 0 {
		t.Fatal("empty wallet identities list in response")
	}

	if walletsListResp.WalletsCount == 0 {
		t.Fatal("empty wallet identities count no equal with expected")
	}

	if walletsListResp.Bookmarks == nil {
		t.Fatal("missing wallet identities bookmarks map in response")
	}

	if len(walletsListResp.Bookmarks) == 0 {
		t.Fatal("empty wallet identities bookmarks map in response")
	}

	for _, walletIdentifier := range createdWallets {
		identifierIdx, isExists := walletsListResp.Bookmarks[walletIdentifier.WalletUUID]
		if !isExists {
			t.Log("wallet not found in enable wallets list")
			t.Fail()
		}

		if uint32(len(walletsListResp.WalletIdentities)) < identifierIdx {
			t.Log("wrong identifier index value in bookmarks")
			t.Fail()
		}

		identifier := walletsListResp.WalletIdentities[identifierIdx]
		if identifier == nil {
			t.Log("nil identifier value in enable wallet resp")
			t.Fail()
		}

		if identifier.WalletUUID != walletIdentifier.WalletUUID {
			t.Log("created wallet uuid not equal with item from enabled wallet list")
			t.Fail()
		}

		if identifier.WalletHash != walletIdentifier.WalletHash {
			t.Log("created wallet hash not equal with item from enabled wallet list")
			t.Fail()
		}
	}

}
