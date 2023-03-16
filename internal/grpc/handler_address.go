/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

package grpc

import (
	"context"
	"github.com/google/uuid"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/forms"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"github.com/crypto-bundle/bc-wallet-common/pkg/tracer"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddress = "GetDerivationAddress"
)

type GetDerivationAddressHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv getAddressMarshallerService
}

// nolint:funlen // fixme
func (h *GetDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	validationForm := &forms.GetDerivationAddressForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletUUID, err := uuid.Parse(validationForm.WalletUUID)
	if err != nil {
		return nil, err
	}

	mnemonicWalletUUID, err := uuid.Parse(validationForm.MnemonicWalletUUID)
	if err != nil {
		return nil, err
	}

	addressData, err := h.walletSrv.GetAddressByPath(ctx, walletUUID, mnemonicWalletUUID,
		validationForm.AccountIndex, validationForm.InternalIndex, validationForm.AddressIndex)
	if err != nil {
		return nil, err
	}

	return &pbApi.DerivationAddressResponse{
		AddressIdentity: &pbApi.DerivationAddressIdentity{
			AccountIndex:  addressData.AccountIndex,
			InternalIndex: addressData.InternalIndex,
			AddressIndex:  addressData.AddressIndex,
			Address:       addressData.Address,
		},
	}, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSrv: walletSrv,
	}
}
