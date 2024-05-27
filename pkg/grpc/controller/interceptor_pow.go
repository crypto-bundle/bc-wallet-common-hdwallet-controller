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
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type powShieldInterceptor struct {
	obscurityDataSvc    obscurityDataProvider
	accessTokensDataSvc accessTokensDataService
	txStatementSvc      transactionalStatementManager

	powTarget *big.Int
}

func (i *powShieldInterceptor) Invoke(ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	switch req.(type) {
	case *GetWalletSessionRequest,
		*GetWalletSessionsRequest,
		*CloseWalletSessionsRequest,
		*GetAccountRequest,
		*GetMultipleAccountRequest,
		*PrepareSignRequestReq,
		*ExecuteSignRequestReq:
		return i.invoke(ctx, method, req, reply, cc, invoker, opts...)
	case *StartWalletSessionRequest:
		return i.invokeSession(ctx, method, req, reply, cc, invoker, opts...)
	case *GetWalletInfoRequest:
		isSystemToken := ctx.Value(app.ContextIsSystemTokenTag).(bool)
		if isSystemToken {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return i.invoke(ctx, method, req, reply, cc, invoker, opts...)

	default:
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *powShieldInterceptor) getObscurityData(ctx context.Context, walletUUID string) ([]byte, error) {
	var obscurityRawData []byte

	lastIdentity, err := i.obscurityDataSvc.GetLastObscurityDataIdentifier(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	if lastIdentity != nil {
		obscurityRawData = lastIdentity[:]

		return obscurityRawData, nil
	}

	accessTokenStr, err := i.accessTokensDataSvc.GetAccessTokenForWallet(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	return []byte(accessTokenStr), nil
}

func (i *powShieldInterceptor) invokeSession(ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	reqRaw, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}

	walletUUID, err := extractWalletUUIDFromReq(req)
	if err != nil {
		return err
	}

	err = i.txStatementSvc.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		obscurityRawData, clbErr := i.getObscurityData(txStmtCtx, walletUUID)

		proofNonce, proofInt := calcPOWProof(i.powTarget, reqRaw, obscurityRawData)
		nonceStr := strconv.FormatInt(proofNonce, 10)

		newCtx := metadata.AppendToOutgoingContext(ctx,
			"X-POW-Hashcash-Proof", fmt.Sprintf("%x", proofInt.Bytes()),
			"X-POW-Hashcash-Nonce", nonceStr)

		clbErr = invoker(newCtx, method, req, reply, cc, opts...)
		if clbErr != nil {
			return clbErr
		}

		resp, isCasted := reply.(*StartWalletSessionResponse)
		if !isCasted {
			return ErrMissingResponse
		}

		clbErr = i.obscurityDataSvc.AddLastObscurityDataIdentifier(txStmtCtx, walletUUID,
			resp.SessionIdentifier.SessionUUID)
		if clbErr != nil {
			return clbErr
		}

		return nil
	})
	if err != nil {
		return err
	}

}

func (i *powShieldInterceptor) invoke(ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	reqRaw, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}

	walletUUID, err := extractWalletUUIDFromReq(req)
	if err != nil {
		return err
	}

	var obscurityRawData []byte

	err = i.txStatementSvc.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		lastIdentity, clbErr := i.obscurityDataSvc.GetLastObscurityDataIdentifier(txStmtCtx)
		if clbErr != nil {
			return clbErr
		}

		if lastIdentity != nil {
			obscurityRawData = lastIdentity[:]

			return nil
		}

		accessTokenStr, clbErr := i.accessTokensDataSvc.GetAccessTokenForWallet(txStmtCtx, walletUUID)
		if clbErr != nil {
			return clbErr
		}

		obscurityRawData = []byte(accessTokenStr)

		return nil
	})
	if err != nil {
		return err
	}

	proofNonce, proofInt := calcPOWProof(i.powTarget, reqRaw, obscurityRawData)
	nonceStr := strconv.FormatInt(proofNonce, 10)

	newCtx := metadata.AppendToOutgoingContext(ctx,
		"X-POW-Hashcash-Proof", fmt.Sprintf("%x", proofInt.Bytes()),
		"X-POW-Hashcash-Nonce", nonceStr)

	return invoker(newCtx, method, req, reply, cc, opts...)
}

func calcPOWProof(powTarget *big.Int,
	protoMsg []byte,
	obscurityData []byte,
) (int64, *big.Int) { // return nonce and proof
	proofNonce := int64(0)
	var reqInt *big.Int = big.NewInt(0)

	for proofNonce != math.MaxInt64 {
		concatRaw := append(protoMsg[:], obscurityData[:]...)
		concatRaw = append(concatRaw, byte(proofNonce))

		hash := sha256.Sum256(concatRaw)
		reqInt.SetBytes(hash[:])

		if reqInt.Cmp(powTarget) == -1 {
			break
		} else {
			proofNonce++
		}
	}

	return proofNonce, reqInt
}

func newPowShieldInterceptor(
	obscurityDataSvc obscurityDataProvider,
	accessTokensDataSvc accessTokensDataService,
	txStatementSvc transactionalStatementManager,
) *powShieldInterceptor {
	powTarget := big.NewInt(1)
	powTarget = powTarget.Lsh(powTarget, uint(256-8))

	return &powShieldInterceptor{
		obscurityDataSvc:    obscurityDataSvc,
		accessTokensDataSvc: accessTokensDataSvc,
		txStatementSvc:      txStatementSvc,

		powTarget: powTarget,
	}
}
