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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller/mocks"
	"go.uber.org/zap/zaptest"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"
)

func TestHdWalletControllerApiClient_GetMultipleAccounts(t *testing.T) {
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

	accessTokenData := make(map[string]string)

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

	accessTokensCount := 5

	createWalletResp, loopErr := managerApiClient.AddNewWallet(ctx, uint(accessTokensCount))
	if loopErr != nil {
		t.Fatal(loopErr)
	}

	if createWalletResp == nil {
		t.Fatal("add new wallet empty response")
	}

	if createWalletResp.WalletIdentifier == nil {
		t.Fatal("missed wallet identifier")
	}

	if createWalletResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_CREATED)
	}

	walletEnableResp, loopErr := managerApiClient.EnableWallet(ctx, createWalletResp.WalletIdentifier.WalletUUID)
	if loopErr != nil {
		t.Fatal(loopErr)
	}

	if walletEnableResp == nil {
		t.Fatal("missing enable wallet response")
	}

	if walletEnableResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_ENABLED {
		t.Fatalf("%s: curent:%s, expected: %s", "wrong wallet status",
			walletEnableResp.WalletStatus, pbCommon.WalletStatus_WALLET_STATUS_ENABLED)
	}

	createRuleToken := createWalletResp.AccessTokens[0]
	accessTokenData[createWalletResp.WalletIdentifier.WalletUUID] = string(createRuleToken.AccessTokenData)

	startWalletSessionResp, loopErr := walletApiClient.StartWalletSession(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID)
	if loopErr != nil {
		t.Fatal(loopErr)
	}

	if startWalletSessionResp == nil {
		t.Fatal("missing start wallet session response")
	}

	addrList := &pbCommon.RangeUnitsList{
		RangeUnits: []*pbCommon.RangeRequestUnit{
			{
				AccountIndex:     5,
				InternalIndex:    10,
				AddressIndexFrom: 106,
				AddressIndexTo:   156,
			},
			{
				AccountIndex:     7,
				InternalIndex:    8,
				AddressIndexFrom: 1201,
				AddressIndexTo:   1251,
			},
			{
				AccountIndex:     10,
				InternalIndex:    32,
				AddressIndexFrom: 401,
				AddressIndexTo:   451,
			},
		},
	}

	addrListAnyData := &anypb.Any{}
	err = addrListAnyData.MarshalFrom(addrList)
	if err != nil {
		t.Fatal(err)
	}

	getAccountResp, err := walletApiClient.GetMultipleAccountsByWallet(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID,
		startWalletSessionResp.SessionIdentifier.SessionUUID,
		addrListAnyData,
	)
	if err != nil {
		t.Fatal(err)
	}

	if getAccountResp == nil {
		t.Fatal("missing get account response")
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

	if getAccountResp.SessionIdentifier.SessionUUID != startWalletSessionResp.SessionIdentifier.SessionUUID {
		t.Fatal("session identifier not equal with expected")
	}

	if getAccountResp.AccountIdentitiesCount == 0 {
		t.Fatal("zero value account identifiers count")
	}

	if getAccountResp.AccountIdentifiers == nil || len(getAccountResp.AccountIdentifiers) == 0 {
		t.Fatal("empty account identifiers list")
	}

	if uint64(len(getAccountResp.AccountIdentifiers)) != getAccountResp.AccountIdentitiesCount {
		t.Fatal("account identifiers count not equal with length of acc identifiers list")
	}

	for _, accountIdentity := range getAccountResp.AccountIdentifiers {
		if len(accountIdentity.Address) == 0 {
			t.Logf("zero account identity address length")
			t.Fail()
		}

		if accountIdentity.Parameters == nil {
			t.Logf("missing account identity parameters")
			t.Fail()
		}
	}

}
