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

package sign_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *Service) PrepareSignRequest(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
	purposeUUID string,
	accountParameters *anypb.Any,
) (addr *pbCommon.AccountIdentity, signReqItem *entities.SignRequest, err error) {
	rawData, err := proto.Marshal(accountParameters)
	if err != nil {
		return nil, nil, err
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		savedItem, clbErr := s.signReqDataSvc.AddSignRequestItem(txStmtCtx, &entities.SignRequest{
			UUID:        uuid.NewString(),
			WalletUUID:  mnemonicUUID,
			SessionUUID: sessionUUID,
			PurposeUUID: purposeUUID,
			AccountData: rawData,
			Status:      types.SignRequestStatusCreated,
			CreatedAt:   time.Time{},
			UpdatedAt:   nil,
		})
		if clbErr != nil {
			return clbErr
		}

		signOwner, clbErr := s.signPrepare(ctx, mnemonicUUID, accountParameters)
		if clbErr != nil {
			return clbErr
		}

		clbErr = s.signReqDataSvc.UpdateSignRequestItemStatus(txStmtCtx, savedItem.UUID,
			types.SignRequestStatusPrepared)
		if clbErr != nil {
			return clbErr
		}

		addr = signOwner
		savedItem.Status = types.SignRequestStatusPrepared
		signReqItem = savedItem

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	err = s.eventPublisherSvc.SendSignPreparedEvent(ctx, signReqItem.UUID)
	if err != nil {
		s.logger.Error("unable to broadcast sign request prepared event", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, mnemonicUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, sessionUUID),
			zap.String(app.SignRequestUUIDTag, signReqItem.UUID))

		// no return - it's ok
	}

	return
}

func (s *Service) signPrepare(ctx context.Context,
	mnemonicUUID string,
	accountParameters *anypb.Any,
) (signerAddr *pbCommon.AccountIdentity, err error) {
	resp, err := s.hdWalletClientSvc.LoadAccount(ctx, &hdwallet.LoadAccountRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: accountParameters,
		},
	})
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(ErrUnableDecodeGrpcErrorStatus))
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
			return nil, nil

		default:
			s.logger.Error("unable to load derivation address from hd-wallet",
				zap.Error(ErrUnableDecodeGrpcErrorStatus))

			return nil, err
		}
	}

	return resp.AccountIdentifier, nil
}
