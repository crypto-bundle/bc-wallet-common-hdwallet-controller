package grpc

import (
	"bc-wallet-eth-hdwallet/internal/grpc/handlers"
	"context"

	"bc-wallet-eth-hdwallet/internal/app"
	"bc-wallet-eth-hdwallet/internal/config"
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hd_wallet_api/proto"

	"go.uber.org/zap"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    *config.Config

	walletSrv walleter

	// all GRPC handlers
	addNewWalletHandler         *handlers.AddNewWalletHandler
	getDerivationAddressHandler *handlers.GetDerivationAddressHandler
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

// NewGRPCHandler New instance of service
// nolint:revive // fixme
func NewGRPCHandler(ctx context.Context,
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

		addNewWalletHandler:         handlers.MakeAddNewWalletHandler(l, walletSrv),
		getDerivationAddressHandler: handlers.MakeGetDerivationAddressHandler(l, walletSrv),
	}, nil
}
