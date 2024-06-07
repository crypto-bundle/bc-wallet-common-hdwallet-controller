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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller/mocks"
	"go.uber.org/zap/zaptest"
	"testing"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
)

func TestHdWalletControllerApiClient_GetEnabledWallets(t *testing.T) {
	logger := zaptest.NewLogger(t)
	clientCfg, err := mocks.NewManagerApiClientConfig("localhost",
		8114, "tron", "./test_case_data/root_token_data.json")
	if err != nil {
		t.Fatal(err)
	}

	client := NewManagerApiClientWrapper(logger, clientCfg.GetRootToken())

	ctx := context.Background()
	err = client.Init(ctx, clientCfg)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Dial(ctx)
	if err != nil {
		t.Fatal(err)
	}

	accessTokensCount := 5

	createdWallets := make(map[string]*pbCommon.MnemonicWalletIdentity, 3)

	for i := 0; i != 3; i++ {
		createWalletResp, loopErr := client.AddNewWallet(ctx, uint(accessTokensCount))
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if createWalletResp == nil {
			t.Fatal("add new wallet empty response")
		}

		walletEnableResp, loopErr := client.EnableWallet(ctx, createWalletResp.WalletIdentifier.WalletUUID)
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if walletEnableResp == nil {
			t.Fatal("missing enable wallet response")
		}

		createdWallets[createWalletResp.WalletIdentifier.WalletUUID] = createWalletResp.WalletIdentifier
	}

	walletsListResp, err := client.GetEnabledWallets(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if walletsListResp == nil {
		t.Fatal("missing wallet list response")
	}

	if walletsListResp.WalletsData == nil {
		t.Fatal("missing wallet identities list in response")
	}

	if len(walletsListResp.WalletsData) == 0 {
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

		if uint32(len(walletsListResp.WalletsData)) < identifierIdx {
			t.Log("wrong identifier index value in bookmarks")
			t.Fail()
		}

		wData := walletsListResp.WalletsData[identifierIdx]
		if wData == nil {
			t.Log("nil identifier value in enable wallet resp")
			t.Fail()
		}

		if wData.WalletIdentifier.WalletUUID != walletIdentifier.WalletUUID {
			t.Log("created wallet uuid not equal with item from enabled wallet list")
			t.Fail()
		}

		if wData.WalletIdentifier.WalletHash != walletIdentifier.WalletHash {
			t.Log("created wallet hash not equal with item from enabled wallet list")
			t.Fail()
		}
	}

}
