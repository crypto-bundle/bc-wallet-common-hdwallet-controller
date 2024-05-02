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
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodSignPrepare = "PrepareSign"
)

type SignPrepareHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService

	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignPrepareHandler) Handle(ctx context.Context,
	req *pbApi.PrepareSignRequestReq,
) (*pbApi.PrepareSignRequestResponse, error) {
	var err error

	vf := &SignRequestPrepareForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, sessionItem, err := h.walletSvc.GetWalletSessionInfo(ctx, vf.WalletUUID, vf.SessionUUID)
	if err != nil {
		h.l.Error("unable get wallet and wallet session info", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	if sessionItem == nil {
		return nil, status.Error(codes.ResourceExhausted, "mnemonic wallet session not found or expired")
	}

	if !sessionItem.IsSessionActive() {
		return nil, status.Error(codes.ResourceExhausted, "mnemonic wallet session not found or expired")
	}

	signOwner, signReqItem, err := h.signManagerSvc.PrepareSignRequest(ctx, vf.WalletUUID, vf.SessionUUID,
		vf.PurposeUUID,
		vf.AccountParameters)
	if err != nil {
		h.l.Error("unable to sign transaction", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, err.Error())
	}

	if signOwner == nil || signReqItem == nil {
		return nil, status.Error(codes.ResourceExhausted,
			"signer account not found or signature session expired")
	}

	return &pbApi.PrepareSignRequestResponse{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentifier: &pbApi.WalletSessionIdentity{
			SessionUUID: sessionItem.UUID,
		},
		AccountIdentifier: signOwner,
		SignatureRequestInfo: &pbApi.SignRequestData{
			Identifier: &pbApi.SignRequestIdentity{UUID: signReqItem.UUID},
			Status:     pbApi.SignRequestData_ReqStatus(signReqItem.Status),
			CreateAt:   uint64(signReqItem.CreatedAt.Unix()),
		},
	}, nil
}

func MakeSignPrepareHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	signManagerSvc signManagerService,
	marshallerSrv marshallerService,
) *SignPrepareHandler {
	return &SignPrepareHandler{
		l:              loggerEntry.With(zap.String(MethodNameTag, MethodSignPrepare)),
		signManagerSvc: signManagerSvc,
		walletSvc:      walletSrv,
		marshallerSrv:  marshallerSrv,
	}
}
