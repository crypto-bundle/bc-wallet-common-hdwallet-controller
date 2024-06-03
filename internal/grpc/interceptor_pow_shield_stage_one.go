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

package grpc

import (
	"context"
	"encoding/hex"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"math/big"
)

type powShieldPreValidationInterceptor struct {
	powValidationSvc powValidatorService
}

func (i *powShieldPreValidationInterceptor) Handle(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	switch req.(type) {
	case *pbApi.StartWalletSessionRequest,
		*pbApi.GetWalletSessionRequest,
		*pbApi.GetWalletSessionsRequest,
		*pbApi.CloseWalletSessionsRequest,
		*pbApi.GetAccountRequest,
		*pbApi.GetMultipleAccountRequest,
		*pbApi.PrepareSignRequestReq,
		*pbApi.ExecuteSignRequestReq:
		return i.handle(ctx, req, info, handler)
	case *pbApi.GetWalletInfoRequest:
		isSystemToken := ctx.Value(app.ContextIsSystemTokenTag).(bool)
		if isSystemToken {
			return handler(ctx, req)
		}

		return i.handle(ctx, req, info, handler)
	default:
		return handler(ctx, req)
	}
}

func (i *powShieldPreValidationInterceptor) handle(ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to retrieve metadata context")
	}

	data := md.Get(PowShieldProofHeader)
	if data == nil {
		return nil, status.Error(codes.InvalidArgument,
			"missing proof of work hashcash header data in metadata")
	}

	if len(data) > 1 {
		return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash data")
	}

	decodedHex, err := hex.DecodeString(data[0])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash data - not hex string")
	}

	powProofHash := big.NewInt(0).SetBytes(decodedHex)
	//if !ok {
	//	return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash value. unable decode sting")
	//}

	isValid := i.powValidationSvc.PreValidate(ctx, powProofHash.Bytes())
	if !isValid {
		return nil, status.Error(codes.InvalidArgument, "pow shield validation failed")
	}

	return handler(ctx, req)
}

func newPowShieldPreValidationInterceptor(powValidationSvc powValidatorService) grpc.UnaryServerInterceptor {
	str := powShieldPreValidationInterceptor{
		powValidationSvc: powValidationSvc,
	}

	return str.Handle
}
