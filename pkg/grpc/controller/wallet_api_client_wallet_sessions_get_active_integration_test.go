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

func TestHdWalletControllerApiClient_GetAllWalletSessions(t *testing.T) {
	const (
		sessionCountToCreate = uint(3)
	)

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

	createdSessions := make([]*WalletSessionIdentity, sessionCountToCreate)
	createdSessionsMap := make(map[string]*WalletSessionIdentity, len(createdSessions))
	createdSessionsFoundMap := make(map[string]bool, len(createdSessions))

	for i := uint(0); i != sessionCountToCreate; i++ {
		startSessionResp, loopErr := walletApiClient.StartWalletSession(ctx,
			walletEnableResp.WalletIdentifier.WalletUUID)
		if loopErr != nil {
			t.Fatal(loopErr)
		}

		if startSessionResp == nil {
			t.Fatal("missing start wallet session response")
		}

		if startSessionResp.WalletIdentifier == nil {
			t.Fatal("missing wallet identifier in start wallet session resp")
		}

		if startSessionResp.SessionIdentifier == nil {
			t.Fatal("missing wallet session identifier in start wallet session resp")
		}

		if startSessionResp.WalletIdentifier.WalletUUID != walletEnableResp.WalletIdentifier.WalletUUID {
			t.Fatal("missing wallet identifier in start wallet session resp")
		}

		createdSessions[i] = startSessionResp.SessionIdentifier
		createdSessionsMap[startSessionResp.SessionIdentifier.SessionUUID] = startSessionResp.SessionIdentifier
		createdSessionsFoundMap[startSessionResp.SessionIdentifier.SessionUUID] = false

		t.Logf("created new wallet session: %d, %s", i,
			startSessionResp.SessionIdentifier.SessionUUID)
	}

	getSessionResp, err := walletApiClient.GetActiveWalletSessions(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal(err)
	}

	if getSessionResp == nil {
		t.Fatal("missing start wallet session response")
	}

	if getSessionResp.WalletIdentifier == nil {
		t.Fatal("missing wallet identifier in get wallet session response")
	}

	if getSessionResp.ActiveSessions == nil || len(getSessionResp.ActiveSessions) == 0 {
		t.Fatal("missing active sessions in get active sessions resp")
	}

	for _, sessionInfo := range getSessionResp.ActiveSessions {
		sessionUUID := sessionInfo.SessionIdentifier.SessionUUID

		createdSessionIdentifier, isExists := createdSessionsMap[sessionUUID]
		if !isExists {
			continue
		}

		_, isExists = createdSessionsFoundMap[sessionUUID]
		if !isExists {
			continue
		}

		if sessionInfo.SessionStatus != WalletSessionStatus_WALLET_SESSION_STATUS_PREPARED {
			t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet session status",
				sessionInfo.SessionStatus,
				WalletSessionStatus_WALLET_SESSION_STATUS_PREPARED)
		}

		if sessionInfo.SessionIdentifier.SessionUUID != createdSessionIdentifier.SessionUUID {
			t.Fatalf("%s: current: %s , expected: %s", "mismatched session UUIDs",
				sessionInfo.SessionIdentifier.SessionUUID, createdSessionIdentifier.SessionUUID)
		}

		createdSessionsFoundMap[sessionUUID] = true
	}

	for _, sessionIdentifier := range createdSessions {
		isFoundInResp, isExists := createdSessionsFoundMap[sessionIdentifier.SessionUUID]
		if !isExists {
			t.Fatalf("%s: uuid: %s", "missed session identifier in found session map",
				sessionIdentifier.SessionUUID)
		}

		if !isFoundInResp {
			t.Fatalf("%s: uuid: %s", "created session not found in active sessions response",
				sessionIdentifier.SessionUUID)
		}
	}

}
