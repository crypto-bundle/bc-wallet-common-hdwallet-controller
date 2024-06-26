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

package grpc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type accessTokenValidationInterceptor struct {
	jwtSvc       jwtService
	tokenDataSvc accessTokenDataService

	readerWalletApprovedRoles []types.AccessTokenRole
	signApprovedRoles         []types.AccessTokenRole
	signPrepareApprovedRoles  []types.AccessTokenRole
}

func (i accessTokenValidationInterceptor) Handle(ctx context.Context,
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
		*pbApi.EnableWalletsRequest:
		return nil, status.Errorf(codes.PermissionDenied, "wallet access token cant call method: %s",
			info.FullMethod)
	default:
		return i.handleWalletRequest(context.WithValue(ctx, app.ContextIsSystemTokenTag, false),
			req, accessToken, accessTokenHash, info, handler)
	}
}

func (i accessTokenValidationInterceptor) handleWalletRequest(ctx context.Context,
	req any,
	accessToken string,
	accessTokenHash string,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	var requiredRole []types.AccessTokenRole

	switch req.(type) {
	case *pbApi.PrepareSignRequestReq:
		requiredRole = i.signPrepareApprovedRoles
	case *pbApi.ExecuteSignRequestReq:
		requiredRole = i.signApprovedRoles
	case *pbApi.GetWalletInfoRequest:
		requiredRole = i.readerWalletApprovedRoles
	default:
		requiredRole = i.readerWalletApprovedRoles
	}

	data, err := i.jwtSvc.GetTokenData(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"unable to extract data from JWT-token: %s", err)
	}

	tokenUUIDStr, isExist := data[app.JWTTokenUUIDLabel]
	if !isExist {
		return nil, status.Errorf(codes.InvalidArgument,
			"missing token_uuid field: %s", ErrMissingTokenUUIDIdentity)
	}

	tokenUUIDRaw, err := uuid.Parse(tokenUUIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"wrong format of token_uuid value: %s", err)
	}

	tokenExpiredAtStr, isExist := data[app.JWTTokenExpiredLabel]
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
		return nil, status.Errorf(codes.PermissionDenied,
			"access token not found")
	}

	if accessTokenItem.Hash != accessTokenHash {
		return nil, status.Error(codes.PermissionDenied, "access token hash mismatched")
	}

	timeDiff := accessTokenItem.ExpiredAt.Sub(expiredAt)
	if timeDiff != 0 {
		return nil, status.Error(codes.PermissionDenied, "expired_at wrong value - not equal with expected")
	}

	if accessTokenItem.ExpiredAt.Before(time.Now()) {
		return nil, status.Error(codes.PermissionDenied, "access token already expired")
	}

	if slices.Index(requiredRole, accessTokenItem.Role) == -1 {
		return nil, status.Error(codes.PermissionDenied, "access token has no permission")
	}

	newCtx := context.WithValue(ctx, app.ContextTokenUUIDTag, accessTokenItem.UUID)

	return handler(newCtx, req)
}

func newAccessTokenInterceptor(jwtSvc jwtService,
	tokenDataSvc accessTokenDataService,
) grpc.UnaryServerInterceptor {
	str := accessTokenValidationInterceptor{
		jwtSvc:       jwtSvc,
		tokenDataSvc: tokenDataSvc,

		readerWalletApprovedRoles: []types.AccessTokenRole{
			types.AccessTokenRoleReader,
			types.AccessTokenRoleSigner,
			types.AccessTokenRoleFakeSigner,
		},
		signApprovedRoles: []types.AccessTokenRole{
			types.AccessTokenRoleSigner,
		},
		signPrepareApprovedRoles: []types.AccessTokenRole{
			types.AccessTokenRoleSigner,
			types.AccessTokenRoleFakeSigner,
		},
	}

	return str.Handle
}
