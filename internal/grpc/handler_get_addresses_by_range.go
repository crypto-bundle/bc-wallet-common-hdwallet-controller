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
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddressByRange = "GetDerivationAddressByRange"
)

type GetDerivationAddressByRangeHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error

	vf := &derivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, sessionItem, err := h.walletSvc.GetWalletSessionInfo(ctx, vf.MnemonicWalletUUID, vf.SessionUUID)
	if err != nil {
		h.l.Error("unable get wallet and wallet session info", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	if walletItem.Status == types.MnemonicWalletStatusDisabled {
		return nil, status.Error(codes.ResourceExhausted, "wallet disabled")
	}

	if sessionItem == nil || !sessionItem.IsSessionActive() {
		return nil, status.Error(codes.ResourceExhausted, "wallet session not found or already expired")
	}

	addressesList, err := h.walletSvc.GetAddressesByRange(ctx, vf.MnemonicWalletUUID,
		vf.Ranges)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	response, err := h.marshallerSrv.MarshallGetAddressByRange(walletItem, addressesList,
		uint64(vf.RangeSize))
	if err != nil {
		h.l.Error("unable to marshall get addresses data", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSvc marshallerService,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSvc:     walletSvc,
		marshallerSrv: marshallerSvc,
	}
}
