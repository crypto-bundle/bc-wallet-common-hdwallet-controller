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

	"github.com/google/uuid"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHdWalletControllerApiClient_EnableWallet(t *testing.T) {
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

	createWalletResp, err := client.AddNewWallet(ctx, &AddNewWalletRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if createWalletResp == nil {
		t.Fatal("add new wallet empty response")
	}

	walletEnableResp, err := client.EnableWallet(ctx, &EnableWalletRequest{
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

	_, err = uuid.Parse(walletEnableResp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal("wrong wallet identity format, not uuid")
	}

	if len(walletEnableResp.WalletIdentifier.WalletHash) != 64 {
		t.Fatal("wrong length of wallet hash string")
	}

	if walletEnableResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_ENABLED {
		t.Fatal("wallet status not equal with expected", pbCommon.WalletStatus_WALLET_STATUS_ENABLED)
	}

	if walletEnableResp.WalletIdentifier.WalletUUID != createWalletResp.WalletIdentifier.WalletUUID {
		t.Fatal("wallet uuid not equal.",
			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
	}

	if walletEnableResp.WalletIdentifier.WalletHash != createWalletResp.WalletIdentifier.WalletHash {
		t.Fatal("wallet hash not equal.",
			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
	}

	walletInfoResp, err := client.GetWalletInfo(ctx, &GetWalletInfoRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if walletInfoResp == nil {
		t.Fatal("missing wallet info resp")
	}

	if walletInfoResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_ENABLED {
		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_ENABLED)
	}
}
