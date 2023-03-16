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
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/forms"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	"github.com/google/uuid"
	"sync"

	"github.com/crypto-bundle/bc-wallet-common/pkg/tracer"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddressByRange = "GetDerivationAddressByRange"
)

type GetDerivationAddressByRangeHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv getAddressByRangeMarshallerService
}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	validationForm := &forms.DerivationAddressByRangeForm{}
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

	walletsData, err := h.walletSrv.GetAddressesByPathByRange(ctx, walletUUID, mnemonicWalletUUID,
		validationForm.AccountIndex, validationForm.InternalIndex,
		validationForm.AddressIndexFrom, validationForm.AddressIndexTo)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	rangeSize := validationForm.AddressIndexTo - validationForm.AddressIndexFrom

	response := &pbApi.DerivationAddressByRangeResponse{
		AddressIdentities: make([]*pbApi.DerivationAddressIdentity, rangeSize+1),
	}

	wg := sync.WaitGroup{}
	wg.Add(int(rangeSize) + 1)
	for i := uint32(0); i != rangeSize; i++ {
		go func(index uint32) {
			response.AddressIdentities[index] = &pbApi.DerivationAddressIdentity{
				AccountIndex:  walletsData[index].AccountIndex,
				InternalIndex: walletsData[index].InternalIndex,
				AddressIndex:  walletsData[index].AddressIndex,
				Address:       walletsData[index].Address,
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSrv: walletSrv,
	}
}
