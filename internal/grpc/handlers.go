package grpc

import (
	"context"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/app"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/config"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/grpc/handlers"
	pbApi "github.com/cryptowize-tech/bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    *config.Config

	walletSrv walleter

	// all GRPC handlers
	addNewWalletHandler                *handlers.AddNewWalletHandler
	getDerivationAddressHandler        *handlers.GetDerivationAddressHandler
	getDerivationAddressByRangeHandler *handlers.GetDerivationAddressByRangeHandler
	getEnabledWalletsHandler           *handlers.GetEnabledWalletsHandler
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
	cfg *config.Config,
	loggerSrv *zap.Logger,

	walletSrv walleter,
) (pbApi.HdWalletApiServer, error) {

	l := loggerSrv.Named("grpc.server.handler").With(
		zap.String(app.ApplicationNameTag, app.ApplicationName),
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		cfg:                            cfg,
		logger:                         l,

		walletSrv: walletSrv,

		addNewWalletHandler:                handlers.MakeAddNewWalletHandler(l, walletSrv),
		getDerivationAddressHandler:        handlers.MakeGetDerivationAddressHandler(l, walletSrv),
		getEnabledWalletsHandler:           handlers.MakeGetEnabledWalletsHandler(l, walletSrv),
		getDerivationAddressByRangeHandler: handlers.MakeGetDerivationAddressByRangeHandler(l, walletSrv),
	}, nil
}
