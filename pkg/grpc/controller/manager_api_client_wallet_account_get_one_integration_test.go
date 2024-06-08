/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package controller

import (
	"context"
	"testing"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller/mocks"

	"github.com/google/uuid"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHdWalletControllerApiClient_GetAccount(t *testing.T) {
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

	resp, err := client.AddNewWallet(ctx, uint(accessTokensCount))
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("add new wallet empty response")
	}

	if resp.WalletIdentifier == nil {
		t.Fatal("missed wallet identifier")
	}

	if resp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_CREATED)
	}

	walletEnableResp, err := client.EnableWallet(ctx, resp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal(err)
	}

	if walletEnableResp == nil {
		t.Fatal("missing enable wallet response")
	}

	if walletEnableResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_ENABLED {
		t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet status",
			walletEnableResp.WalletStatus, pbCommon.WalletStatus_WALLET_STATUS_ENABLED)
	}

	addr := &pbCommon.DerivationAddressIdentity{
		AccountIndex:  1009,
		InternalIndex: 8,
		AddressIndex:  7,
	}

	anyData := &anypb.Any{}
	err = anyData.MarshalFrom(addr)
	if err != nil {
		t.Fatal(err)
	}

	getAccountResp, err := client.GetWalletAccount(ctx, walletEnableResp.WalletIdentifier.WalletUUID,
		anyData)
	if err != nil {
		t.Fatal(err)
	}

	if getAccountResp == nil || getAccountResp.AccountIdentifier == nil {
		t.Fatal("missing get account response or missing account identifier")
	}

	if getAccountResp.WalletIdentifier == nil {
		t.Fatal("missing wallet identifier in get account resp")
	}

	if getAccountResp.SessionIdentifier == nil {
		t.Fatal("missing session identifier in get account resp")
	}

	if getAccountResp.WalletIdentifier.WalletUUID != walletEnableResp.WalletIdentifier.WalletUUID {
		t.Fatal("wallet identifier not equal with expected")
	}

	if getAccountResp.SessionIdentifier.SessionUUID != uuid.Nil.String() {
		t.Fatal("session identifier not equal with expected and must be zero")
	}
}
