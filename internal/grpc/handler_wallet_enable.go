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

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameEnableWallet = "EnableWallet"
)

type EnableWalletHandler struct {
	l *zap.Logger

	walletSvc walletManagerService

	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *EnableWalletHandler) Handle(ctx context.Context,
	req *pbApi.EnableWalletRequest,
) (*pbApi.EnableWalletResponse, error) {
	var err error

	validationForm := &EnableWalletForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	wallet, err := h.walletSvc.EnableWalletByUUID(ctx, validationForm.WalletUUID)
	if err != nil {
		h.l.Error("unable to enable mnemonic wallet", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if wallet == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	return &pbApi.EnableWalletResponse{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: wallet.UUID.String(),
			WalletHash: wallet.MnemonicHash,
		},
	}, nil
}

func MakeEnableWalletHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *EnableWalletHandler {
	return &EnableWalletHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodNameEnableWallet)),
		walletSvc: walletSvc,
	}
}