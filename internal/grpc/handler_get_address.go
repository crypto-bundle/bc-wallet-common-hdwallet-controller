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
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"sync"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddress = "GetDerivationAddress"
)

type GetDerivationAddressHandler struct {
	l *zap.Logger

	walletSvc     walletManagerService
	marshallerSvc marshallerService

	pbAddrPool *sync.Pool
}

// nolint:funlen // fixme
func (h *GetDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error

	vf := &GetDerivationAddressForm{}
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

	addressData, err := h.walletSvc.GetAddress(ctx, vf.MnemonicWalletUUID,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex)
	if err != nil {
		return nil, err
	}

	if addressData == nil {
		return nil, status.Error(codes.ResourceExhausted,
			"wallet not found or all wallet sessions already expired")
	}

	addressEntity := h.pbAddrPool.Get().(*pbCommon.DerivationAddressIdentity)
	addressEntity.AccountIndex = vf.AccountIndex
	addressEntity.InternalIndex = vf.InternalIndex
	addressEntity.AddressIndex = vf.AddressIndex
	addressEntity.Address = *addressData

	marshalledData, err := h.marshallerSvc.MarshallGetAddressData(walletItem, addressEntity)
	if err != nil {
		h.l.Error("unable to marshall public address data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer func() {
		h.pbAddrPool.Put(addressEntity)
	}()

	return marshalledData, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSrv marshallerService,
	pbAddrPool *sync.Pool,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSvc:     walletSvc,
		marshallerSvc: marshallerSrv,
		pbAddrPool:    pbAddrPool,
	}
}
