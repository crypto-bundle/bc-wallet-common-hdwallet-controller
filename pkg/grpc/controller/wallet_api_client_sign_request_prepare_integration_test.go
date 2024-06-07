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
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHdWalletControllerApiClient_PrepareSignRequest(t *testing.T) {
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

	createWalletResp, err := client.AddNewWallet(ctx, uint(accessTokensCount))
	if err != nil {
		t.Fatal(err)
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

	walletEnableResp, err := client.EnableWallet(ctx, createWalletResp.WalletIdentifier.WalletUUID)
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

	startWalletSessionResp, err := walletApiClient.StartWalletSession(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal(err)
	}

	if startWalletSessionResp == nil {
		t.Fatal("missing start wallet session response")
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

	preparedSignReqResp, err := walletApiClient.PrepareSignRequest(ctx,
		walletEnableResp.WalletIdentifier.WalletUUID,
		startWalletSessionResp.SessionIdentifier.SessionUUID,
		uuid.NewString(),
		anyData)
	if err != nil {
		t.Fatal(err)
	}

	if preparedSignReqResp == nil {
		t.Fatal("missing prepare signReq response")
	}

	if preparedSignReqResp.WalletIdentifier == nil {
		t.Fatal("missing wallet identifier in prepare signReq response")
	}

	if preparedSignReqResp.SessionIdentifier == nil {
		t.Fatal("missing session identifier in prepare signReq response")
	}

	if preparedSignReqResp.AccountIdentifier == nil {
		t.Fatal("missing account identifier in prepare signReq response")
	}

	if preparedSignReqResp.WalletIdentifier.WalletUUID != walletEnableResp.WalletIdentifier.WalletUUID {
		t.Fatal("wallet identifier not equal with expected")
	}

	if preparedSignReqResp.SessionIdentifier.SessionUUID != startWalletSessionResp.SessionIdentifier.SessionUUID {
		t.Fatal("session identifier not equal with expected")
	}

	if preparedSignReqResp.SignatureRequestInfo == nil {
		t.Fatal("missing signReq info")
	}

	if preparedSignReqResp.SignatureRequestInfo.Identifier == nil {
		t.Fatal("missing signReq identifier")
	}

	_, err = uuid.Parse(preparedSignReqResp.SignatureRequestInfo.Identifier.UUID)
	if err != nil {
		t.Fatal("wrong signReq identity format, not uuid")
	}

	if preparedSignReqResp.SignatureRequestInfo.Status != SignRequestData_REQUEST_PREPARED {
		t.Fatalf("%s: curent:%s, expected: %s", "wrong signReq status",
			preparedSignReqResp.SignatureRequestInfo.Status, SignRequestData_REQUEST_PREPARED)
	}
}
