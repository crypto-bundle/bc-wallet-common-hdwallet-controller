package grpc

import (
	"context"
	"sync"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/config"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    *config.MangerConfig

	walletSrv     walletManagerService
	marshallerSrv marshallerService
	// all GRPC handlers
	addNewWalletHandler                *AddNewWalletHandler
	getDerivationAddressHandler        *GetDerivationAddressHandler
	getDerivationAddressByRangeHandler *GetDerivationAddressByRangeHandler
	getEnabledWalletsHandler           *GetEnabledWalletsHandler
	getWalletInfoHandler               *GetWalletInfoHandler
	signTransactionHandle              *SignTransactionHandler
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

func (h *grpcServerHandle) SignTransaction(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	return h.signTransactionHandle.Handle(ctx, req)
}

func (h *grpcServerHandle) GetWalletInfo(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	return h.getWalletInfoHandler.Handle(ctx, req)
}

// New instance of service
func New(ctx context.Context,
	loggerSrv *zap.Logger,

	walletSrv walletManagerService,
) pbApi.HdWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler").With(
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	addrRespPool := &sync.Pool{New: func() any {
		return new(pbApi.DerivationAddressIdentity)
	}}

	marshallerSrv := newGRPCMarshaller(loggerSrv, addrRespPool)

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,

		addNewWalletHandler:         MakeAddNewWalletHandler(l, walletSrv, marshallerSrv),
		getDerivationAddressHandler: MakeGetDerivationAddressHandler(l, walletSrv, marshallerSrv, addrRespPool),
		getEnabledWalletsHandler:    MakeGetEnabledWalletsHandler(l, walletSrv, marshallerSrv),
		getDerivationAddressByRangeHandler: MakeGetDerivationAddressByRangeHandler(l, walletSrv,
			marshallerSrv, addrRespPool),
		signTransactionHandle: MakeSignTransactionsHandler(l, walletSrv, marshallerSrv),
		getWalletInfoHandler:  MakeGetWalletInfoHandler(l, walletSrv, marshallerSrv),
	}
}
