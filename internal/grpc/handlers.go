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

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/config"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletControllerApiServer

	logger *zap.Logger
	cfg    *config.MangerConfig

	powValidatorSvc powValidatorService

	marshallerSrv marshallerService
	// all GRPC handlers
	addNewWalletHandlerSvc   addWalletHandlerService
	importWalletHandlerSvc   importWalletHandlerService
	disableWalletHandlerSvc  disableWalletHandlerService
	disableWalletsHandlerSvc disableWalletsHandlerService
	enableWalletsHandlerSvc  enableWalletsHandlerService
	enableWalletHandlerSvc   enableWalletHandlerService
	getWalletInfoHandlerSvc  getWalletHandlerService
	getEnabledWalletsHandler getEnabledWalletsHandlerService

	startWalletSessionHandlerSvc startWalletSessionHandlerService
	closeWalletSessionHandlerSvc closeWalletSessionHandlerService
	getWalletSessionHandleSvc    getWalletSessionHandlerService
	getWalletSessionsHandleSvc   getWalletSessionsHandlerService

	getAccountHandlerSvc          getAccountHandlerService
	getMultipleAccountsHandlerSvc getAccountsHandlerService

	prepareSignReqHandlerSvc prepareSignRequestHandlerService
	executeSignReqHandleSvc  executeSignRequestHandlerService
}

// Wallet management handler

func (h *grpcServerHandle) AddNewWallet(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	return h.addNewWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) EnableWallet(ctx context.Context,
	req *pbApi.EnableWalletRequest,
) (*pbApi.EnableWalletResponse, error) {
	return h.enableWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) ImportWallet(ctx context.Context,
	req *pbApi.ImportWalletRequest,
) (*pbApi.ImportWalletResponse, error) {
	return h.importWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetWalletInfo(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	return h.getWalletInfoHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetEnabledWallets(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	return h.getEnabledWalletsHandler.Handle(ctx, req)
}

func (h *grpcServerHandle) DisableWallet(ctx context.Context,
	req *pbApi.DisableWalletRequest,
) (*pbApi.DisableWalletResponse, error) {
	return h.disableWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) DisableWallets(ctx context.Context,
	req *pbApi.DisableWalletsRequest,
) (*pbApi.DisableWalletsResponse, error) {
	return h.disableWalletsHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) EnableWallets(ctx context.Context,
	req *pbApi.EnableWalletsRequest,
) (*pbApi.EnableWalletsResponse, error) {
	return h.enableWalletsHandlerSvc.Handle(ctx, req)
}

// Wallet sessions handlers

func (h *grpcServerHandle) StartWalletSession(ctx context.Context,
	req *pbApi.StartWalletSessionRequest,
) (*pbApi.StartWalletSessionResponse, error) {
	return h.startWalletSessionHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetWalletSession(ctx context.Context,
	req *pbApi.GetWalletSessionRequest,
) (*pbApi.GetWalletSessionResponse, error) {
	return h.getWalletSessionHandleSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetAllWalletSessions(ctx context.Context,
	req *pbApi.GetWalletSessionsRequest,
) (*pbApi.GetWalletSessionsResponse, error) {
	return h.getWalletSessionsHandleSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) CloseWalletSession(ctx context.Context,
	req *pbApi.CloseWalletSessionsRequest,
) (*pbApi.CloseWalletSessionsResponse, error) {
	return h.closeWalletSessionHandlerSvc.Handle(ctx, req)
}

// Wallet derivation address handlers

func (h *grpcServerHandle) GetAccount(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	return h.getAccountHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetMultipleAccounts(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	return h.getMultipleAccountsHandlerSvc.Handle(ctx, req)
}

// Sign flow handlers

func (h *grpcServerHandle) PrepareSignRequest(ctx context.Context,
	req *pbApi.PrepareSignRequestReq,
) (*pbApi.PrepareSignRequestResponse, error) {
	return h.prepareSignReqHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) ExecuteSignRequest(ctx context.Context,
	req *pbApi.ExecuteSignRequestReq,
) (*pbApi.ExecuteSignRequestResponse, error) {
	return h.executeSignReqHandleSvc.Handle(ctx, req)
}

// New instance of service
func New(loggerSrv *zap.Logger,
	walletManagerSvc walletManagerService,
	signManagerSvc signManagerService,
) pbApi.HdWalletControllerApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbCommon.DerivationAddressIdentity)
	}}

	marshallerSvc := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerHandle{
		UnimplementedHdWalletControllerApiServer: &pbApi.UnimplementedHdWalletControllerApiServer{},
		logger:                                   l,

		marshallerSrv: marshallerSvc,
		// handlers
		addNewWalletHandlerSvc:   MakeAddNewWalletHandler(l, walletManagerSvc, marshallerSvc),
		importWalletHandlerSvc:   MakeImportWalletHandler(l, walletManagerSvc),
		enableWalletHandlerSvc:   MakeEnableWalletHandler(l, walletManagerSvc),
		getWalletInfoHandlerSvc:  MakeGetWalletInfoHandler(l, walletManagerSvc),
		getEnabledWalletsHandler: MakeGetEnabledWalletsHandler(l, walletManagerSvc, marshallerSvc),
		disableWalletHandlerSvc:  MakeDisableWalletHandler(l, walletManagerSvc, signManagerSvc),
		disableWalletsHandlerSvc: MakeDisableWalletsHandler(l, walletManagerSvc, signManagerSvc),
		enableWalletsHandlerSvc:  MakeEnableWalletsHandler(l, walletManagerSvc),

		startWalletSessionHandlerSvc: MakeStartWalletSessionHandler(l, walletManagerSvc),
		getWalletSessionHandleSvc:    MakeGetWalletSessionHandler(l, walletManagerSvc),
		getWalletSessionsHandleSvc:   MakeGetWalletSessionsHandler(l, walletManagerSvc, marshallerSvc),
		closeWalletSessionHandlerSvc: MakeCloseWalletSessionHandler(l, walletManagerSvc, signManagerSvc),

		getAccountHandlerSvc: MakeGetAccountHandler(l, walletManagerSvc,
			marshallerSvc, addrRespPool),
		getMultipleAccountsHandlerSvc: MakeGetMultipleAccountsHandler(l, walletManagerSvc,
			marshallerSvc),
		prepareSignReqHandlerSvc: MakeSignPrepareHandler(l, walletManagerSvc, signManagerSvc, marshallerSvc),
		executeSignReqHandleSvc:  MakeSignTransactionsHandler(l, walletManagerSvc, signManagerSvc, marshallerSvc),
	}
}
