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
	"sync"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetAccountSessionLessHandler struct {
	l *zap.Logger

	walletSvc     walletManagerService
	marshallerSvc marshallerService

	pbAddrPool *sync.Pool
}

// nolint:funlen // fixme
func (h *GetAccountSessionLessHandler) Handle(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	var err error

	vf := &GetAccountSessionLessForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, err := h.walletSvc.GetWalletByUUID(ctx, vf.MnemonicWalletUUID)
	if err != nil {
		h.l.Error("unable get wallet and wallet session info", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	if walletItem.Status == types.MnemonicWalletStatusDisabled {
		return nil, status.Error(codes.ResourceExhausted, "wallet disabled")
	}

	addressData, err := h.walletSvc.GetAccountWithoutSession(ctx, walletItem,
		vf.AccountParameters)
	if err != nil {
		return nil, err
	}

	if addressData == nil {
		return nil, status.Error(codes.ResourceExhausted,
			"wallet not found or all wallet sessions already expired")
	}

	return &pbApi.GetAccountResponse{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentifier: &pbApi.WalletSessionIdentity{
			SessionUUID: uuid.Nil.String(),
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: vf.AccountParameters,
			Address:    *addressData,
		},
	}, nil
}

func MakeGetAccountSessionLessHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSrv marshallerService,
	pbAddrPool *sync.Pool,
) *GetAccountSessionLessHandler {
	return &GetAccountSessionLessHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetAccount)),
		walletSvc:     walletSvc,
		marshallerSvc: marshallerSrv,
		pbAddrPool:    pbAddrPool,
	}
}
