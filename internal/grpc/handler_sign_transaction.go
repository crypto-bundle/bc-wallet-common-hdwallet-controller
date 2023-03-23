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

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodSignTransaction = "SignTransaction"
)

type SignTransactionHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignTransactionHandler) Handle(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	walletUUID, err := uuid.Parse(req.WalletUUID)
	if err != nil {
		return nil, err
	}

	mnemonicWalletUUID, err := uuid.Parse(req.MnemonicWalletUUID)
	if err != nil {
		return nil, err
	}

	signedTxData, err := h.walletSrv.SignTransaction(ctx, walletUUID, mnemonicWalletUUID,
		req.AddressIdentity.AccountIndex, req.AddressIdentity.InternalIndex, req.AddressIdentity.AddressIndex,
		&tronCore.Transaction{})
	if err != nil {
		h.l.Error("unable to sign transaction", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := h.marshallerSrv.MarshallSignTransaction(signedTxData)
	if err != nil {
		h.l.Error("unable to marshall signed transaction data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func MakeSignTransactionsHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *SignTransactionHandler {
	return &SignTransactionHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodSignTransaction)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
