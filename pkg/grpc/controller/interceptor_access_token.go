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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type accessTokenInterceptor struct {
	accessTokensDataSvc accessTokensDataService
}

func (i *accessTokenInterceptor) Invoke(ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	switch req.(type) {
	case *AddNewWalletRequest, *ImportWalletRequest:
		return invoker(ctx, method, req, reply, cc, opts...)
	default:
		return i.invoke(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func (i *accessTokenInterceptor) invoke(ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	var walletUUID string
	switch req.(type) {
	case *GetWalletInfoRequest:
		walletUUID = req.(*GetWalletInfoRequest).WalletIdentifier.WalletUUID
	case *EnableWalletsRequest:
		walletUUID = req.(*EnableWalletRequest).WalletIdentifier.WalletUUID
	case *StartWalletSessionRequest:
		walletUUID = req.(*StartWalletSessionRequest).WalletIdentifier.WalletUUID

	default:
		return i.invoke(ctx, method, req, reply, cc, invoker, opts...)
	}

	accessTokenStr, err := i.accessTokensDataSvc.GetAccessTokenForWallet(ctx, walletUUID)
	if err != nil {
		return err
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(make(map[string]string))
	}

	md.Set("X-Access-Token", accessTokenStr)

	return i.invoke(metadata.NewOutgoingContext(ctx, md),
		method, req, reply, cc, invoker, opts...)
}
