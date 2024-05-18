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
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/protobuf/proto"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/google/uuid"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestHdWalletControllerApiClient_GetWalletInfo_SystemToken(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		//originGRPC.WithChainUnaryInterceptor(),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	systemTokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2YWx1ZXNfbWFwIjp7Imluc3RhbGxtZW50X3V1aWQiOiI5MzBkZTBjMy1mZDVmLTRmYTQtYmNhYS1jZWQwNWE2OGVhZGUiLCJleHBpcmVkX2F0IjoiMjA3NC0wNS0wNiAxNjoxNjowNyJ9fQ.yonSdLs9cZTq2ER4-4y-4kWMHx85J0eRqQFRCXkqhOY"

	newCtx := metadata.AppendToOutgoingContext(ctx, "X-Access-Token", systemTokenStr)

	createWalletResp, err := client.AddNewWallet(newCtx, &AddNewWalletRequest{
		CreateAccessTokensCount: 5,
	})
	if err != nil {
		t.Fatal(err)
	}

	if createWalletResp == nil {
		t.Fatal("add new wallet empty response")
	}

	//tokenStr := string(createWalletResp.AccessTokens[0].AccessTokenData)

	req := &GetWalletInfoRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
		},
	}

	//reqRaw, err := proto.Marshal(req)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//target := big.NewInt(1)
	//target.Lsh(target, uint(256-8))
	//
	//target2 := big.NewInt(1)
	//target2.Lsh(target2, uint(256-24))
	//
	//nonce := int64(0)
	//
	////t.Logf("target: %d", target)
	////t.Logf("target2: %d", target2)
	//
	//var reqInt *big.Int = big.NewInt(0)
	//
	//for nonce != math.MaxInt64 {
	//	concatRaw := append(reqRaw[:], byte(nonce))
	//
	//	hash := sha256.Sum256(concatRaw)
	//	reqInt.SetBytes(hash[:])
	//
	//	t.Logf("\r%x: calc: %d, target: %d", hash, reqInt, target)
	//
	//	if reqInt.Cmp(target) == -1 {
	//		break
	//	} else {
	//		nonce++
	//	}
	//}
	//
	//t.Logf("target str: %x", target.String())
	//t.Logf("pow    str: %x", reqInt.String())

	//md := metadata.Pairs("X-Access-Token", systemTokenStr,
	//	"X-POW-Hashcash-Proof", reqInt.String(),
	//	"X-POW-Hashcash-Nonce", strconv.FormatInt(nonce, 10),
	//)

	//newCtx := metadata.NewOutgoingContext(ctx, md)

	walletInfoResp, err := client.GetWalletInfo(newCtx, req)
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

func TestHdWalletControllerApiClient_GetWalletInfo_WalletToken(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		//originGRPC.WithChainUnaryInterceptor(),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	systemTokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2YWx1ZXNfbWFwIjp7Imluc3RhbGxtZW50X3V1aWQiOiI5MzBkZTBjMy1mZDVmLTRmYTQtYmNhYS1jZWQwNWE2OGVhZGUiLCJleHBpcmVkX2F0IjoiMjA3NC0wNS0wNiAxNjoxNjowNyJ9fQ.yonSdLs9cZTq2ER4-4y-4kWMHx85J0eRqQFRCXkqhOY"

	newCtx := metadata.AppendToOutgoingContext(ctx, "X-Access-Token", systemTokenStr)

	createWalletResp, err := client.AddNewWallet(newCtx, &AddNewWalletRequest{
		CreateAccessTokensCount: 5,
	})
	if err != nil {
		t.Fatal(err)
	}

	if createWalletResp == nil {
		t.Fatal("add new wallet empty response")
	}

	signerTokenStr := string(createWalletResp.AccessTokens[0].AccessTokenData)
	accessTokenUUIDStr := createWalletResp.AccessTokens[0].AccessTokenIdentifier.UUID

	accessTokenUUID, err := uuid.Parse(accessTokenUUIDStr)
	if err != nil {
		t.Fatal("wrong access token uuid format, not uuid")
	}

	req := &GetWalletInfoRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: createWalletResp.WalletIdentifier.WalletUUID,
		},
	}

	reqRaw, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(256-4))
	nonce := int64(0)

	//t.Logf("target: %d", target)
	//t.Logf("target2: %d", target2)

	var reqInt *big.Int = big.NewInt(0)
	var hash []byte

	h := sha256.New()

	for nonce != math.MaxInt64 {
		concatRaw := append(reqRaw[:], accessTokenUUID[:]...)
		concatRaw = append(concatRaw, byte(nonce))

		h.Reset()
		h.Write(concatRaw)
		hash = h.Sum(nil)

		reqInt.SetBytes(bytes.Clone(hash))

		t.Logf("\r%x: calc: %d, target: %d", hash, reqInt, target)

		if reqInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	t.Logf("target str: %x", target.String())
	t.Logf("pow    str: %x", reqInt.String())
	tragetHash := fmt.Sprintf("%x", reqInt.Bytes())

	md := metadata.Pairs("X-Access-Token", signerTokenStr,
		"X-POW-Hashcash-Proof", tragetHash,
		"X-POW-Hashcash-Nonce", strconv.FormatInt(nonce, 10),
	)

	newCtx = metadata.NewOutgoingContext(ctx, md)

	walletInfoResp, err := client.GetWalletInfo(newCtx, req)
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
