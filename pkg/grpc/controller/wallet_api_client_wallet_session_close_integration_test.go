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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller/mocks"
	"go.uber.org/zap/zaptest"
	"testing"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
)

func TestHdWalletControllerApiClient_CloseWalletSession(t *testing.T) {
	logger := zaptest.NewLogger(t)
	clientCfg, err := mocks.NewManagerApiClientConfig("localhost",
		8114, "tron", "./test_case_data/root_token_data.json")
	if err != nil {
		t.Fatal(err)
	}

	managerApiClient := NewManagerApiClientWrapper(logger, clientCfg.GetRootToken())

	ctx := context.Background()
	err = managerApiClient.Init(ctx, clientCfg)
	if err != nil {
		t.Fatal(err)
	}

	err = managerApiClient.Dial(ctx)
	if err != nil {
		t.Fatal(err)
	}

	accessTokensCount := 5

	createWalletResp, err := managerApiClient.AddNewWallet(ctx, uint(accessTokensCount))
	if err != nil {
		t.Fatal(err)
	}

	if createWalletResp == nil {
		t.Fatal("add new wallet empty response")
	}

	createRuleToken := createWalletResp.AccessTokens[0]

	accessTokenData := make(map[string]string, 1)
	accessTokenData[createWalletResp.WalletIdentifier.WalletUUID] = string(createRuleToken.AccessTokenData)

	walletApiConfig := mocks.NewWalletApiClientConfig("localhost",
		8115, "tron")

	walletApiClient := NewWalletApiClientWrapper(logger,
		mocks.NewObscurityDataStoreStore(make(map[string][]byte)),
		mocks.NewAccessTokenDataStore(accessTokenData),
		mocks.NewTxStmtMock())

	err = walletApiClient.Init(ctx, walletApiConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = walletApiClient.Dial(ctx)
	if err != nil {
		t.Fatal(err)
	}

	walletEnableResp, err := managerApiClient.EnableWallet(ctx, createWalletResp.WalletIdentifier.WalletUUID)
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

	startWalletSessionResp, loopErr := walletApiClient.StartWalletSession(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID)
	if loopErr != nil {
		t.Fatal(loopErr)
	}

	if startWalletSessionResp == nil {
		t.Fatal("missing start wallet session response")
	}

	if startWalletSessionResp.WalletIdentifier == nil {
		t.Fatal("missing wallet identifier in start wallet session resp")
	}

	if startWalletSessionResp.SessionIdentifier == nil {
		t.Fatal("missing wallet session identifier in start wallet session resp")
	}

	if startWalletSessionResp.WalletIdentifier.WalletUUID != walletEnableResp.WalletIdentifier.WalletUUID {
		t.Fatal("missing wallet identifier in start wallet session resp")
	}

	closeWalletSessionResp, err := walletApiClient.CloseWalletSession(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID,
		startWalletSessionResp.SessionIdentifier.SessionUUID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if closeWalletSessionResp == nil {
		t.Fatal("missing start wallet session response")
	}

	if closeWalletSessionResp.WalletIdentifier == nil {
		t.Fatal("missing wallet identifier in get wallet session response")
	}

	if closeWalletSessionResp.SessionIdentifier == nil {
		t.Fatal("missing session identifier in get wallet session response")
	}

	if closeWalletSessionResp.SessionStatus != WalletSessionStatus_WALLET_SESSION_STATUS_CLOSED {
		t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet session status",
			closeWalletSessionResp.SessionStatus,
			WalletSessionStatus_WALLET_SESSION_STATUS_CLOSED)
	}
}
