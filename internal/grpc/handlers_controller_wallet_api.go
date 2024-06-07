/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
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

// grpcServerWalletApiHandler is wrapper struct for implementation all grpc handlers
type grpcServerWalletApiHandler struct {
	*pbApi.UnimplementedHdWalletControllerWalletApiServer

	logger *zap.Logger
	cfg    *config.MangerConfig

	powValidatorSvc powValidatorService

	marshallerSrv marshallerService
	// all GRPC handlers
	getWalletInfoHandlerSvc getWalletHandlerService

	startWalletSessionHandlerSvc startWalletSessionHandlerService
	closeWalletSessionHandlerSvc closeWalletSessionHandlerService
	getWalletSessionHandleSvc    getWalletSessionHandlerService
	getWalletSessionsHandleSvc   getWalletSessionsHandlerService

	getAccountHandlerSvc          getAccountHandlerService
	getMultipleAccountsHandlerSvc getAccountsHandlerService

	prepareSignReqHandlerSvc prepareSignRequestHandlerService
	executeSignReqHandleSvc  executeSignRequestHandlerService
}

func (h *grpcServerWalletApiHandler) GetWalletInfo(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	return h.getWalletInfoHandlerSvc.Handle(ctx, req)
}

// Wallet sessions handlers

func (h *grpcServerWalletApiHandler) StartWalletSession(ctx context.Context,
	req *pbApi.StartWalletSessionRequest,
) (*pbApi.StartWalletSessionResponse, error) {
	return h.startWalletSessionHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerWalletApiHandler) GetWalletSession(ctx context.Context,
	req *pbApi.GetWalletSessionRequest,
) (*pbApi.GetWalletSessionResponse, error) {
	return h.getWalletSessionHandleSvc.Handle(ctx, req)
}

func (h *grpcServerWalletApiHandler) GetAllWalletSessions(ctx context.Context,
	req *pbApi.GetWalletSessionsRequest,
) (*pbApi.GetWalletSessionsResponse, error) {
	return h.getWalletSessionsHandleSvc.Handle(ctx, req)
}

func (h *grpcServerWalletApiHandler) CloseWalletSession(ctx context.Context,
	req *pbApi.CloseWalletSessionsRequest,
) (*pbApi.CloseWalletSessionsResponse, error) {
	return h.closeWalletSessionHandlerSvc.Handle(ctx, req)
}

// Wallet Sub-account address handlers

func (h *grpcServerWalletApiHandler) GetAccount(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	return h.getAccountHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerWalletApiHandler) GetMultipleAccounts(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	return h.getMultipleAccountsHandlerSvc.Handle(ctx, req)
}

// Sign flow handlers

func (h *grpcServerWalletApiHandler) PrepareSignRequest(ctx context.Context,
	req *pbApi.PrepareSignRequestReq,
) (*pbApi.PrepareSignRequestResponse, error) {
	return h.prepareSignReqHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerWalletApiHandler) ExecuteSignRequest(ctx context.Context,
	req *pbApi.ExecuteSignRequestReq,
) (*pbApi.ExecuteSignRequestResponse, error) {
	return h.executeSignReqHandleSvc.Handle(ctx, req)
}

// NewWalletApiHandler instance of service
func NewWalletApiHandler(loggerSrv *zap.Logger,
	walletManagerSvc walletManagerService,
	signManagerSvc signManagerService,
) pbApi.HdWalletControllerWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbCommon.DerivationAddressIdentity)
	}}

	marshallerSvc := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerWalletApiHandler{
		UnimplementedHdWalletControllerWalletApiServer: &pbApi.UnimplementedHdWalletControllerWalletApiServer{},

		logger: l,

		marshallerSrv: marshallerSvc,
		// handlers
		getWalletInfoHandlerSvc: MakeGetWalletInfoHandler(l, walletManagerSvc),

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
