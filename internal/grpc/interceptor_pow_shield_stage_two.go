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
	"go.uber.org/zap"
	"math/big"
	"strconv"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type powShieldFullValidationInterceptor struct {
	logger *zap.Logger

	walletDataSvc   walletDataService
	powProofDataSvc powProofDataService

	powValidationSvc powValidatorService

	txStmtSvc transactionalStatementManager
}

func (i *powShieldFullValidationInterceptor) Handle(ctx context.Context,
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

func (i *powShieldFullValidationInterceptor) handle(ctx context.Context,
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
		return nil, status.Error(codes.InvalidArgument, "missing hashcash data in metadata")
	}

	if len(data) > 1 {
		return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash data")
	}

	nonceData := md.Get(PowShieldNonceHeader)
	if data == nil {
		return nil, status.Error(codes.InvalidArgument,
			"missing proof of work hashcash-nonce header data in metadata")
	}

	if len(nonceData) > 1 {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash-nonce data")
	}

	accessTokenUUID, ok := ctx.Value(app.ContextTokenUUIDTag).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash data")
	}

	nonce, err := strconv.ParseInt(nonceData[0], 10, 0)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash-nonce data")
	}

	return i.validateForSessionFlow(ctx, data[0], nonce,
		req, accessTokenUUID, handler)
}

func (i *powShieldFullValidationInterceptor) validateForSessionFlow(ctx context.Context,
	hashHexString string,
	nonce int64,
	req any,
	accessTokenUUID uuid.UUID,
	unaryHandler grpc.UnaryHandler,
) (resp any, err error) {
	protoMsg, ok := req.(proto.Message)
	if !ok {
		return nil, err
	}

	decodedHex, err := hex.DecodeString(hashHexString)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash data - not hex string")
	}

	powProofHash := big.NewInt(0).SetBytes(decodedHex)

	err = i.txStmtSvc.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		powProof, clbErr := i.powProofDataSvc.GetPowProofByMessageHash(txStmtCtx,
			hashHexString)
		if clbErr != nil {
			i.logger.Error("unable to check pow exist", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		if powProof != nil {
			return status.Error(codes.InvalidArgument, "pow hash already used")
		}

		lastSessionIdentity, clbErr := i.walletDataSvc.GetLastWalletSessionIdentityByAccessTokenUUID(txStmtCtx,
			accessTokenUUID.String())
		if clbErr != nil {
			i.logger.Error("unable to get last wallet session identity", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		obscurityData := accessTokenUUID
		if lastSessionIdentity != nil {
			obscurityData = lastSessionIdentity.SessionUUID
		}

		protoMsgRawData, clbErr := proto.Marshal(protoMsg)
		if clbErr != nil {
			i.logger.Error("unable to marshal body message", zap.Error(clbErr))

			return clbErr
		}

		isProofDataValid, clbErr := i.powValidationSvc.ValidateByObscurityData(txStmtCtx, powProofHash.Bytes(), nonce,
			protoMsgRawData, obscurityData)
		if clbErr != nil {
			i.logger.Error("unable to validate pow", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		if !isProofDataValid {
			i.logger.Warn("pow validation failed",
				zap.String(app.PowHashTag, hashHexString))

			return status.Error(codes.InvalidArgument, "pow shield validation failed")
		}

		currentTime := time.Now()
		_, clbErr = i.powProofDataSvc.AddNewPowProof(txStmtCtx, &entities.PowProof{
			UUID:            uuid.NewString(),
			AccessTokenUUID: accessTokenUUID.String(),

			MessageCheckNonce: nonce,
			MessageHash:       hashHexString,
			MessageData:       protoMsgRawData,

			CreatedAt: time.Now(),
			UpdatedAt: &currentTime,
		})
		if clbErr != nil {
			i.logger.Error("unable to save new pow item", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		handlerResp, clbErr := unaryHandler(ctx, req)
		if clbErr != nil { // in this case clbErr already content status.Error
			return clbErr
		}

		resp = handlerResp

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newPowShieldFullValidationInterceptor(logger *zap.Logger,
	walletSessionDataSvc walletDataService,
	powProofDataSvc powProofDataService,
	powValidationSvc powValidatorService,
	txStmtSvc transactionalStatementManager,
) grpc.UnaryServerInterceptor {
	str := powShieldFullValidationInterceptor{
		logger:           logger,
		walletDataSvc:    walletSessionDataSvc,
		powProofDataSvc:  powProofDataSvc,
		powValidationSvc: powValidationSvc,
		txStmtSvc:        txStmtSvc,
	}

	return str.Handle
}
