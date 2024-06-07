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
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller/mocks"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/google/uuid"
)

func TestHdWalletWalletApiClient_GetWalletInfo(t *testing.T) {
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

	walletInfoResp, err := walletApiClient.GetWalletInfo(ctx, createWalletResp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal(err)
	}

	if walletInfoResp == nil {
		t.Fatal("missing wallet info resp")
	}

	_, err = uuid.Parse(walletInfoResp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal("wrong wallet identity format, not uuid")
	}

	if len(walletInfoResp.WalletIdentifier.WalletHash) != 64 {
		t.Fatal("wrong length of wallet hash string")
	}

	if walletInfoResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_CREATED)
	}

	if walletInfoResp.WalletIdentifier.WalletUUID != createWalletResp.WalletIdentifier.WalletUUID {
		t.Fatal("wallet uuid not equal.",
			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
	}

	if walletInfoResp.WalletIdentifier.WalletHash != createWalletResp.WalletIdentifier.WalletHash {
		t.Fatal("wallet hash not equal.",
			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
	}
}

//func TestHdWalletControllerApiClient_GetWalletInfo_WalletToken(t *testing.T) {
//	options := []originGRPC.DialOption{
//		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
//		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
//		originGRPC.WithBlock(),
//		//originGRPC.WithChainUnaryInterceptor(),
//	}
//	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
//	if err != nil {
//		t.Error(err)
//	}
//
//	client := NewHdWalletControllerApiClient(grpcConn)
//	ctx := context.Background()
//
//	systemTokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2YWx1ZXNfbWFwIjp7Imluc3RhbGxtZW50X3V1aWQiOiI5MzBkZTBjMy1mZDVmLTRmYTQtYmNhYS1jZWQwNWE2OGVhZGUiLCJleHBpcmVkX2F0IjoiMjA3NC0wNS0wNiAxNjoxNjowNyJ9fQ.yonSdLs9cZTq2ER4-4y-4kWMHx85J0eRqQFRCXkqhOY"
//
//	newCtx := metadata.AppendToOutgoingContext(ctx, "X-Access-Token", systemTokenStr)
//
//	createWalletResp, err := client.AddNewWallet(newCtx, &AddNewWalletRequest{
//		CreateAccessTokensCount: 5,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if createWalletResp == nil {
//		t.Fatal("add new wallet empty response")
//	}
//
//	signerTokenStr := string(createWalletResp.AccessTokens[0].AccessTokenData)
//	accessTokenUUIDStr := createWalletResp.AccessTokens[0].AccessTokenIdentifier.UUID
//
//	accessTokenUUID, err := uuid.Parse(accessTokenUUIDStr)
//	if err != nil {
//		t.Fatal("wrong access token uuid format, not uuid")
//	}
//
//	req := &GetWalletInfoRequest{
//		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
//			WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
//		},
//	}
//
//	reqRaw, err := proto.Marshal(req)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	target := big.NewInt(1)
//	target.Lsh(target, uint(256-4))
//	nonce := int64(0)
//
//	//t.Logf("target: %d", target)
//	//t.Logf("target2: %d", target2)
//
//	var reqInt *big.Int = big.NewInt(0)
//	var hash []byte
//
//	h := sha256.New()
//
//	for nonce != math.MaxInt64 {
//		concatRaw := append(reqRaw[:], accessTokenUUID[:]...)
//		concatRaw = append(concatRaw, byte(nonce))
//
//		h.Reset()
//		h.Write(concatRaw)
//		hash = h.Sum(nil)
//
//		reqInt.SetBytes(bytes.Clone(hash))
//
//		t.Logf("\r%x: calc: %d, target: %d", hash, reqInt, target)
//
//		if reqInt.Cmp(target) == -1 {
//			break
//		} else {
//			nonce++
//		}
//	}
//
//	t.Logf("target str: %x", target.String())
//	t.Logf("pow    str: %x", reqInt.String())
//	tragetHash := fmt.Sprintf("%x", reqInt.Bytes())
//
//	md := metadata.Pairs("X-Access-Token", signerTokenStr,
//		"X-POW-Hashcash-Proof", tragetHash,
//		"X-POW-Hashcash-Nonce", strconv.FormatInt(nonce, 10),
//	)
//
//	newCtx = metadata.NewOutgoingContext(ctx, md)
//
//	walletInfoResp, err := client.GetWalletInfo(newCtx, req)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if walletInfoResp == nil {
//		t.Fatal("missing wallet info resp")
//	}
//
//	_, err = uuid.Parse(walletInfoResp.WalletIdentifier.WalletUUID)
//	if err != nil {
//		t.Fatal("wrong wallet identity format, not uuid")
//	}
//
//	if len(walletInfoResp.WalletIdentifier.WalletHash) != 64 {
//		t.Fatal("wrong length of wallet hash string")
//	}
//
//	if walletInfoResp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
//		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_CREATED)
//	}
//
//	if walletInfoResp.WalletIdentifier.WalletUUID != createWalletResp.WalletIdentifier.WalletUUID {
//		t.Fatal("wallet uuid not equal.",
//			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
//	}
//
//	if walletInfoResp.WalletIdentifier.WalletHash != createWalletResp.WalletIdentifier.WalletHash {
//		t.Fatal("wallet hash not equal.",
//			"wallet identifier of 'create wallet' and 'get wallet info' requests not equal")
//	}
//}
