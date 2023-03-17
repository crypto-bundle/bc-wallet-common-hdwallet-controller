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
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/config"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    *config.Config

	walletSrv     walletManagerService
	marshallerSrv marshallerService
	// all GRPC handlers
	addNewWalletHandler                *AddNewWalletHandler
	getDerivationAddressHandler        *GetDerivationAddressHandler
	getDerivationAddressByRangeHandler *GetDerivationAddressByRangeHandler
	getEnabledWalletsHandler           *GetEnabledWalletsHandler
}

func (h *grpcServerHandle) AddNewWallet(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	return h.addNewWalletHandler.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddress(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	return h.getDerivationAddressHandler.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddressByRange(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	return h.getDerivationAddressByRangeHandler.Handle(ctx, req)
}

func (h *grpcServerHandle) GetEnabledWallets(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	return h.getEnabledWalletsHandler.Handle(ctx, req)
}

// New instance of service
func New(ctx context.Context,
	loggerSrv *zap.Logger,

	walletSrv walletManagerService,
) (pbApi.HdWalletApiServer, error) {

	l := loggerSrv.Named("grpc.server.handler").With(
		zap.String(app.ApplicationNameTag, app.ApplicationName),
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	marshallerSrv := newGRPCMarshaller(loggerSrv)

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		walletSrv: walletSrv,

		marshallerSrv:                      marshallerSrv,
		addNewWalletHandler:                MakeAddNewWalletHandler(l, walletSrv, marshallerSrv),
		getDerivationAddressHandler:        MakeGetDerivationAddressHandler(l, walletSrv, marshallerSrv),
		getEnabledWalletsHandler:           MakeGetEnabledWalletsHandler(l, walletSrv, marshallerSrv),
		getDerivationAddressByRangeHandler: MakeGetDerivationAddressByRangeHandler(l, walletSrv, marshallerSrv),
	}, nil
}
