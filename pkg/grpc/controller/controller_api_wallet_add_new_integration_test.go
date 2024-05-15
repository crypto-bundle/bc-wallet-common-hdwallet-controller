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
	"github.com/crypto-bundle/bc-wallet-common-lib-jwt/pkg/jwt"
	"testing"
	"time"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/google/uuid"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHdWalletControllerApiClient_AddNewWallet(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		//originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		//originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	jwtSvc := jwt.NewJWTService("123456")
	expiredTime := time.Now().Add(time.Hour * 24 * 356 * 7)
	claim := jwt.NewTokenClaimBuilder(expiredTime)
	tokenUUID := uuid.NewString()
	_ = claim.AddData("token_uuid", tokenUUID)
	_ = claim.AddData("token_expired_at", expiredTime.Format(time.DateTime))
	tokenStr, _ := jwtSvc.GenerateJWT(claim)

	resp, err := client.AddNewWallet(ctx, &AddNewWalletRequest{
		AccessTokens: []*AccessTokenData{
			{
				AccessTokenIdentifier: &AccessTokenIdentity{UUID: tokenUUID},
				AccessTokenData:       []byte(tokenStr),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("add new wallet empty response")
	}

	if resp.WalletIdentifier == nil {
		t.Fatal("missed wallet identifier")
	}

	_, err = uuid.Parse(resp.WalletIdentifier.WalletUUID)
	if err != nil {
		t.Fatal("wrong wallet identity format, not uuid")
	}

	if len(resp.WalletIdentifier.WalletHash) != 64 {
		t.Fatal("wrong length of wallet hash string")
	}

	if resp.WalletStatus != pbCommon.WalletStatus_WALLET_STATUS_CREATED {
		t.Fatal("wallet status not equal with expected:", pbCommon.WalletStatus_WALLET_STATUS_CREATED)
	}
}
