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

// grpcServerMangerApiHandle is wrapper struct for implementation all grpc handlers
type grpcServerMangerApiHandle struct {
	*pbApi.UnimplementedHdWalletControllerManagerApiServer

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

	getAccountHandlerSvc getAccountHandlerService
}

// Wallet management handler

func (h *grpcServerMangerApiHandle) AddNewWallet(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	return h.addNewWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) EnableWallet(ctx context.Context,
	req *pbApi.EnableWalletRequest,
) (*pbApi.EnableWalletResponse, error) {
	return h.enableWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) ImportWallet(ctx context.Context,
	req *pbApi.ImportWalletRequest,
) (*pbApi.ImportWalletResponse, error) {
	return h.importWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) GetWalletInfo(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	return h.getWalletInfoHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) GetEnabledWallets(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	return h.getEnabledWalletsHandler.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) DisableWallet(ctx context.Context,
	req *pbApi.DisableWalletRequest,
) (*pbApi.DisableWalletResponse, error) {
	return h.disableWalletHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) DisableWallets(ctx context.Context,
	req *pbApi.DisableWalletsRequest,
) (*pbApi.DisableWalletsResponse, error) {
	return h.disableWalletsHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerMangerApiHandle) EnableWallets(ctx context.Context,
	req *pbApi.EnableWalletsRequest,
) (*pbApi.EnableWalletsResponse, error) {
	return h.enableWalletsHandlerSvc.Handle(ctx, req)
}

// Wallet Sub-account address handler

func (h *grpcServerMangerApiHandle) GetAccount(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	return h.getAccountHandlerSvc.Handle(ctx, req)
}

// NewManagerApiHandler instance of service
func NewManagerApiHandler(loggerSrv *zap.Logger,
	walletManagerSvc walletManagerService,
	signManagerSvc signManagerService,
) pbApi.HdWalletControllerManagerApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbCommon.DerivationAddressIdentity)
	}}

	marshallerSvc := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerMangerApiHandle{
		UnimplementedHdWalletControllerManagerApiServer: &pbApi.UnimplementedHdWalletControllerManagerApiServer{},

		logger: l,

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

		getAccountHandlerSvc: MakeGetAccountSessionLessHandler(l, walletManagerSvc,
			marshallerSvc, addrRespPool),
	}
}
