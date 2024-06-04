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
	"crypto/sha256"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type accessRootTokenValidationInterceptor struct {
	systemTokenHash string

	jwtSvc       jwtService
	tokenDataSvc accessTokenDataService

	readerWalletApprovedRoles []types.AccessTokenRole
	signApprovedRoles         []types.AccessTokenRole
	signPrepareApprovedRoles  []types.AccessTokenRole
}

func (i accessRootTokenValidationInterceptor) Handle(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
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

	accessToken := accessTokenData[0]
	accessTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(accessToken)))

	switch req.(type) {
	case *pbApi.AddNewWalletRequest,
		*pbApi.ImportWalletRequest,
		*pbApi.EnableWalletRequest,
		*pbApi.GetEnabledWalletsRequest,
		*pbApi.DisableWalletRequest,
		*pbApi.DisableWalletsRequest,
		*pbApi.EnableWalletsRequest,
		*pbApi.GetWalletInfoRequest,
		*pbApi.GetAccountRequest:
		return i.handleSystemRequest(ctx, req, accessToken, accessTokenHash, info, handler)
	default:
		return nil, status.Errorf(codes.PermissionDenied, "root token cant call method: %s",
			info.FullMethod)
	}
}

func (i accessRootTokenValidationInterceptor) handleSystemRequest(ctx context.Context,
	req any,
	_ string, // access token string
	accessTokenHash string,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	if accessTokenHash != i.systemTokenHash {
		return nil, status.Error(codes.PermissionDenied, "wrong token hash value")
	}

	return handler(context.WithValue(ctx, app.ContextIsSystemTokenTag, true), req)
}

func newAccessRootTokenInterceptor(systemTokenHash string,
) grpc.UnaryServerInterceptor {
	str := accessRootTokenValidationInterceptor{
		systemTokenHash: systemTokenHash,
	}

	return str.Handle
}
