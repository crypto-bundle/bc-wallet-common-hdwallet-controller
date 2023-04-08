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
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"sync"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/forms"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

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
	marshallerSrv marshallerService
	respPool      sync.Pool
}

type respAddrList []*pbApi.DerivationAddressIdentity

func (l *respAddrList) Reset() {

}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &forms.DerivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletPubData, err := h.walletSrv.GetWalletByUUID(ctx, vf.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletPubData == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	mnemoWalletData, isExists := walletPubData.MnemonicWalletsByUUID[vf.MnemonicWalletUUIDRaw]
	if !isExists {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	return h.processRequest(ctx, vf, walletPubData, mnemoWalletData)
}

func (h *GetDerivationAddressByRangeHandler) processRequest(ctx context.Context,
	vf *forms.DerivationAddressByRangeForm,
	walletPubData *types.PublicWalletData,
	mnemoWalletData *types.PublicMnemonicWalletData,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error

	rangeSize := (vf.AddressIndexTo - vf.AddressIndexFrom) + 1
	filedData := make([]*pbApi.DerivationAddressIdentity, rangeSize)

	marshallerCallback := func(addressIdx, position uint32, address string) {
		addressEntity := h.respPool.Get().(*pbApi.DerivationAddressIdentity)

		addressEntity.AccountIndex = vf.AccountIndex
		addressEntity.InternalIndex = vf.AccountIndex
		addressEntity.AddressIndex = addressIdx
		addressEntity.Address = address

		filedData[position] = addressEntity
		return
	}

	err = h.walletSrv.GetAddressesByPathByRange(ctx, vf.WalletUUIDRaw, vf.MnemonicWalletUUIDRaw,
		vf.AccountIndex, vf.InternalIndex,
		vf.AddressIndexFrom, vf.AddressIndexTo, marshallerCallback)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	response, err := h.marshallerSrv.MarshallGetAddressByRange(walletPubData, mnemoWalletData, filedData)
	if err != nil {
		h.l.Error("unable to marshall get addresses data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer func(clearedSize uint32) {
		go func(size uint32) {
			for i := uint32(0); i != size; i++ {
				h.respPool.Put(filedData[i])
			}
		}(clearedSize)
	}(rangeSize)

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
		respPool: sync.Pool{New: func() any {
			return new(pbApi.DerivationAddressIdentity)
		}},
	}
}
