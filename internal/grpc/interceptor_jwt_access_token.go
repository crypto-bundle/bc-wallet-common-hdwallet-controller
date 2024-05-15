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
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"github.com/google/uuid"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const (
	AccessTokenHeader = "X-Access-Token"

	ContextTokenUUIDTag = "access_token_uuid"
)

type accessTokenValidationInterceptor struct {
	jwtSvc       jwtService
	tokenDataSvc accessTokenDataService
}

func (i accessTokenValidationInterceptor) Handle(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	switch req.(type) {
	case *pbApi.AddNewWalletRequest, *pbApi.ImportWalletRequest:
		return handler(ctx, req)
	default:
		return i.handle(ctx, req, info, handler)
	}
}

func (i accessTokenValidationInterceptor) handle(ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to retrieve metadata context")
	}

	accessTokenData := md.Get(AccessTokenHeader)
	if accessTokenData == nil {
		return nil, status.Error(codes.InvalidArgument, "missing access token in metadata")
	}

	if len(accessTokenData) > 1 {
		return nil, status.Error(codes.InvalidArgument, "wrong format of access token")
	}

	data, err := i.jwtSvc.GetTokenData(accessTokenData[0])
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"unable to extract data from JWT-token: %s", err)
	}

	tokenUUIDStr, isExist := data[TokenUUIDLabel]
	if !isExist {
		return nil, status.Errorf(codes.InvalidArgument,
			"missing token_uuid field: %s", ErrMissingTokenUUIDIdentity)
	}

	tokenUUIDRaw, err := uuid.Parse(tokenUUIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"wrong format of token_uuid value: %s", err)
	}

	tokenExpiredAtStr, isExist := data[TokenExpiredLabel]
	if !isExist {
		return nil, status.Error(codes.InvalidArgument, "missing expired_at field")
	}

	expiredAt, err := time.Parse(time.DateTime, tokenExpiredAtStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"wrong format of expired_at field, required: %s", time.DateTime)
	}

	accessTokenItem, err := i.tokenDataSvc.GetAccessTokenInfoByUUID(ctx, tokenUUIDRaw.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get access token from store")
	}

	if accessTokenItem == nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"access token not found")
	}

	timeDiff := accessTokenItem.ExpiredAt.Sub(expiredAt)
	if timeDiff != 0 {
		return nil, status.Error(codes.InvalidArgument, "expired_at wrong value - not equal with expected")
	}

	if accessTokenItem.ExpiredAt.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "access token already expired")
	}

	return handler(context.WithValue(ctx, ContextTokenUUIDTag, accessTokenItem.UUID),
		req)
}

func newAccessTokenInterceptor(jwtSvc jwtService,
	tokenDataSvc accessTokenDataService,
) grpc.UnaryServerInterceptor {
	str := accessTokenValidationInterceptor{
		jwtSvc:       jwtSvc,
		tokenDataSvc: tokenDataSvc,
	}

	return str.Handle
}
