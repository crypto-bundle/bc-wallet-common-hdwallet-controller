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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"go.uber.org/zap"
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
	case *pbApi.AddNewWalletRequest, *pbApi.ImportWalletRequest:
		return handler(ctx, req)
	default:
		return i.handle(ctx, req, info, handler)
	}
}

func (i *powShieldFullValidationInterceptor) handle(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to retrieve metadata context")
	}

	data := md.Get(PowShieldTokenHeader)
	if data == nil {
		return nil, status.Error(codes.InvalidArgument, "missing hashcash data in metadata")
	}

	if len(data) > 1 {
		return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash data")
	}

	accessTokenUUID, ok := ctx.Value(ContextTokenUUIDTag).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "wrong format of hashcash data")
	}

	return i.validateForSessionFlow(ctx, []byte(data[0]), req, accessTokenUUID, handler)
}

func (i *powShieldFullValidationInterceptor) validateForSessionFlow(ctx context.Context,
	originHash []byte,
	req any,
	accessTokenUUID uuid.UUID,
	unaryHandler grpc.UnaryHandler,
) (resp any, err error) {
	protoMsg, ok := req.(proto.Message)
	if !ok {
		return nil, err
	}

	err = i.txStmtSvc.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		powProof, clbErr := i.powProofDataSvc.GetPowProofByMessageHash(txStmtCtx,
			originHash)
		if clbErr != nil {
			i.logger.Error("unable to check pow exist", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		if powProof != nil {
			return status.Error(codes.InvalidArgument, "pow hash already used")
		}

		lastSessionIdentities, clbErr := i.walletDataSvc.GetLastWalletSessionIdentityByAccessTokenUUID(txStmtCtx,
			accessTokenUUID.String())
		if clbErr != nil {
			i.logger.Error("unable to get last wallet session identity", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		protoMsgRawData, clbErr := proto.Marshal(protoMsg)
		if clbErr != nil {
			i.logger.Error("unable to marshal body message", zap.Error(clbErr))

			return clbErr
		}

		isProofDataValid, nonce, clbErr := i.powValidationSvc.FullValidate(txStmtCtx, originHash, protoMsgRawData,
			lastSessionIdentities.SessionUUID)
		if clbErr != nil {
			i.logger.Error("unable to validate pow", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		if !isProofDataValid {
			i.logger.Warn("pow validation failed",
				zap.String(app.PowHashTag, string(originHash)))

			return status.Error(codes.InvalidArgument, "pow shield validation failed")
		}

		currentTime := time.Now()
		_, clbErr = i.powProofDataSvc.AddNewPowProof(txStmtCtx, &entities.PowProof{
			UUID:            uuid.NewString(),
			AccessTokenUUID: lastSessionIdentities.AccessTokeUUID.String(),

			MessageCheckNonce: nonce,
			MessageHash:       originHash,
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
	return powShieldFullValidationInterceptor{
		logger:           logger,
		walletDataSvc:    walletSessionDataSvc,
		powProofDataSvc:  powProofDataSvc,
		powValidationSvc: powValidationSvc,
		txStmtSvc:        txStmtSvc,
	}.Handle
}
