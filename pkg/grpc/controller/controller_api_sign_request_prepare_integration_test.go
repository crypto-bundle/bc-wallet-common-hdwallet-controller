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
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHdWalletControllerApiClient_PrepareSignRequest(t *testing.T) {
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

	resp, err := client.AddNewWallet(ctx, &AddNewWalletRequest{})
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

	walletEnableResp, err := client.EnableWallet(ctx, &EnableWalletRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: resp.WalletIdentifier.WalletUUID,
		},
	})
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

	startWalletSessionResp, err := client.StartWalletSession(ctx, &StartWalletSessionRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletEnableResp.WalletIdentifier.WalletUUID,
		},
	})
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

	preparedSignReqResp, err := client.PrepareSignRequest(ctx, &PrepareSignRequestReq{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletEnableResp.WalletIdentifier.WalletUUID,
		},
		SessionIdentifier: &WalletSessionIdentity{
			SessionUUID: startWalletSessionResp.SessionIdentifier.SessionUUID,
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: anyData,
		},
		SignPurposeIdentifier: &SignPurposeIdentity{
			UUID: uuid.NewString(),
		},
	})
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
