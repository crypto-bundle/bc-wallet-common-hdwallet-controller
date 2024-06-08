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
		*pbApi.GetWalletInfoRequest,
		*pbApi.GetAccountRequest,
		*pbApi.GetMultipleAccountRequest,
		*pbApi.PrepareSignRequestReq,
		*pbApi.ExecuteSignRequestReq:
		return i.handle(ctx, req, info, handler)
	default:
		return nil, status.Errorf(codes.PermissionDenied, "token cant call method: %s",
			info.FullMethod)
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
	if nonceData == nil {
		return nil, status.Error(codes.InvalidArgument,
			"missing proof of work hashcash-nonce header data in metadata")
	}

	if len(nonceData) > 1 {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash-nonce data")
	}

	accessTokenData := md.Get(AccessTokenHeader)
	if accessTokenData == nil {
		return nil, status.Error(codes.InvalidArgument, "missing access token in metadata")
	}

	if len(accessTokenData) > 1 {
		return nil, status.Error(codes.InvalidArgument, "wrong format of access token")
	}

	accessToken := accessTokenData[0]

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
		req, accessToken, accessTokenUUID, handler)
}

func (i *powShieldFullValidationInterceptor) validateForSessionFlow(ctx context.Context,
	powHashHexString string,
	powNonce int64,
	req any,
	accessToken string,
	accessTokenUUID uuid.UUID,
	unaryHandler grpc.UnaryHandler,
) (resp any, err error) {
	protoMsg, ok := req.(proto.Message)
	if !ok {
		return nil, err
	}

	decodedHex, err := hex.DecodeString(powHashHexString)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,
			"wrong format of hashcash data - not hex string")
	}

	powProofHash := big.NewInt(0).SetBytes(decodedHex)

	err = i.txStmtSvc.BeginTxWithRollbackOnError(context.Background(), func(txStmtCtx context.Context) error {
		powProof, clbErr := i.powProofDataSvc.GetPowProofByMessageHash(txStmtCtx,
			powHashHexString)
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

		obscurityData := []byte(accessToken)
		if lastSessionIdentity != nil {
			//sessionUUID = lastSessionIdentity.SessionUUID.String()
			obscurityData = lastSessionIdentity.SessionUUID[:]
		}

		//i.logger.Info("data",
		//	zap.String("access_token_uuid", accessTokenUUID.String()),
		//	zap.String("access_token_hash", fmt.Sprintf("%x", sha256.Sum256([]byte(accessToken)))),
		//	zap.String("obsc_data", string(obscurityData)),
		//	zap.String("session_uuid", sessionUUID),
		//)

		protoMsgRawData, clbErr := proto.Marshal(protoMsg)
		if clbErr != nil {
			i.logger.Error("unable to marshal body message", zap.Error(clbErr))

			return clbErr
		}

		isProofDataValid, clbErr := i.powValidationSvc.ValidateByObscurityData(txStmtCtx, powProofHash.Bytes(), powNonce,
			protoMsgRawData, obscurityData)
		if clbErr != nil {
			i.logger.Error("unable to validate pow", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		if !isProofDataValid {
			i.logger.Warn("pow validation failed",
				zap.String(app.PowHashTag, powHashHexString))

			return status.Error(codes.InvalidArgument, "pow shield validation failed")
		}

		currentTime := time.Now()
		_, clbErr = i.powProofDataSvc.AddNewPowProof(txStmtCtx, &entities.PowProof{
			UUID:            uuid.NewString(),
			AccessTokenUUID: accessTokenUUID.String(),

			MessageCheckNonce: powNonce,
			MessageHash:       powHashHexString,
			MessageData:       protoMsgRawData,

			CreatedAt: currentTime,
			UpdatedAt: &currentTime,
		})
		if clbErr != nil {
			i.logger.Error("unable to save new pow item", zap.Error(clbErr))

			return status.Error(codes.Internal, "something went wrong")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	//handlerResp, err := unaryHandler(ctx, req)
	//if err != nil { // in this case clbErr already content status.Error
	//	return nil, err
	//}
	//
	//resp = handlerResp

	return unaryHandler(ctx, req)
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
