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
	disableWalletHandlerSvc  disableWalletHandlerService
	disableWalletsHandlerSvc disableWalletsHandlerService
	enableWalletHandlerSvc   enableWalletHandlerService
	getWalletInfoHandlerSvc  getWalletHandlerService
	getEnabledWalletsHandler getEnabledWalletsHandlerService

	startWalletSessionHandlerSvc startWalletSessionHandlerService
	getWalletSessionHandleSvc    getWalletSessionHandlerService
	getWalletSessionsHandleSvc   getWalletSessionsHandlerService

	getDerivationAddressHandlerSvc        getAddressHandlerService
	getDerivationAddressByRangeHandlerSvc getAddressByRangeHandlerService
	signTransactionHandle                 signTransactionRequestHandlerService
}

func (h *grpcServerHandle) AddNewWallet(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	return h.addNewWalletHandlerSvc.Handle(ctx, req)
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

func (h *grpcServerHandle) GetEnabledWallets(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	return h.getEnabledWalletsHandler.Handle(ctx, req)
}

func (h *grpcServerHandle) SignTransaction(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	return h.signTransactionHandle.Handle(ctx, req)
}

func (h *grpcServerHandle) GetWalletInfo(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	return h.getWalletInfoHandlerSvc.Handle(ctx, req)
}

// New instance of service
func New(ctx context.Context,
	loggerSrv *zap.Logger,

	mnemonicWalletDataSvc mnemonicWalletsDataService,
	walletSessionDataSvc walletSessionDataService,
) pbApi.HdWalletManagerApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbCommon.DerivationAddressIdentity)
	}}

	marshallerSrv := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerHandle{
		UnimplementedHdWalletManagerApiServer: &pbApi.UnimplementedHdWalletManagerApiServer{},
		logger:                                l,

		marshallerSrv: marshallerSrv,
		// handlers
		addNewWalletHandlerSvc:   MakeAddNewWalletHandler(l, marshallerSrv),
		disableWalletHandlerSvc:  nil,
		disableWalletsHandlerSvc: nil,
		enableWalletHandlerSvc:   nil,
		getWalletInfoHandlerSvc:  MakeGetWalletInfoHandler(l, marshallerSrv),

		startWalletSessionHandlerSvc: nil,
		getWalletSessionHandleSvc:    nil,
		getWalletSessionsHandleSvc:   nil,

		getDerivationAddressHandlerSvc: MakeGetDerivationAddressHandler(l, marshallerSrv, addrRespPool),
		getEnabledWalletsHandler:       MakeGetEnabledWalletsHandler(l, marshallerSrv),
		getDerivationAddressByRangeHandlerSvc: MakeGetDerivationAddressByRangeHandler(l,
			marshallerSrv, addrRespPool),
		signTransactionHandle: MakeSignTransactionsHandler(l, marshallerSrv),
	}
}
