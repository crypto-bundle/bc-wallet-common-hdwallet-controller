package grpc

import (
	"context"
	"sync"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/config"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletManagerApiServer

	logger *zap.Logger
	cfg    *config.MangerConfig

	mnemonicWalletDataSvc mnemonicWalletsDataService
	walletSessionDataSvc  walletSessionDataService

	marshallerSrv marshallerService
	// all GRPC handlers
	addNewWalletHandlerSvc   addWalletHandlerService
	importWalletHandlerSvc   importWalletHandlerService
	disableWalletHandlerSvc  disableWalletHandlerService
	disableWalletsHandlerSvc disableWalletsHandlerService
	enableWalletHandlerSvc   enableWalletHandlerService
	getWalletInfoHandlerSvc  getWalletHandlerService
	getEnabledWalletsHandler getEnabledWalletsHandlerService

	startWalletSessionHandlerSvc startWalletSessionHandlerService
	closeWalletSessionHandlerSvc closeWalletSessionHandlerService
	getWalletSessionHandleSvc    getWalletSessionHandlerService
	getWalletSessionsHandleSvc   getWalletSessionsHandlerService

	getDerivationAddressHandlerSvc        getAddressHandlerService
	getDerivationAddressByRangeHandlerSvc getAddressByRangeHandlerService

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

func (h *grpcServerHandle) GetDerivationAddress(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	return h.getDerivationAddressHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddressByRange(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	return h.getDerivationAddressByRangeHandlerSvc.Handle(ctx, req)
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
) pbApi.HdWalletManagerApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbCommon.DerivationAddressIdentity)
	}}

	marshallerSvc := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerHandle{
		UnimplementedHdWalletManagerApiServer: &pbApi.UnimplementedHdWalletManagerApiServer{},
		logger:                                l,

		marshallerSrv: marshallerSvc,
		// handlers
		addNewWalletHandlerSvc:   MakeAddNewWalletHandler(l, walletManagerSvc, marshallerSvc),
		importWalletHandlerSvc:   MakeImportWalletHandler(l, walletManagerSvc),
		enableWalletHandlerSvc:   MakeEnableWalletHandler(l, walletManagerSvc),
		getWalletInfoHandlerSvc:  MakeGetWalletInfoHandler(l, walletManagerSvc),
		getEnabledWalletsHandler: MakeGetEnabledWalletsHandler(l, walletManagerSvc, marshallerSvc),
		disableWalletHandlerSvc:  MakeDisableWalletHandler(l, walletManagerSvc, signManagerSvc),
		disableWalletsHandlerSvc: MakeDisableWalletsHandler(l, walletManagerSvc, signManagerSvc),

		startWalletSessionHandlerSvc: MakeStartWalletSessionHandler(l, walletManagerSvc),
		getWalletSessionHandleSvc:    MakeGetWalletSessionHandler(l, walletManagerSvc),
		getWalletSessionsHandleSvc:   MakeGetWalletSessionsHandler(l, walletManagerSvc, marshallerSvc),
		closeWalletSessionHandlerSvc: MakeCloseWalletSessionHandler(l, walletManagerSvc, signManagerSvc),

		getDerivationAddressHandlerSvc: MakeGetDerivationAddressHandler(l, walletManagerSvc,
			marshallerSvc, addrRespPool),
		getDerivationAddressByRangeHandlerSvc: MakeGetDerivationAddressByRangeHandler(l, walletManagerSvc,
			marshallerSvc),
		prepareSignReqHandlerSvc: MakeSignPrepareHandler(l, walletManagerSvc, signManagerSvc, marshallerSvc),
		executeSignReqHandleSvc:  MakeSignTransactionsHandler(l, walletManagerSvc, signManagerSvc, marshallerSvc),
	}
}
