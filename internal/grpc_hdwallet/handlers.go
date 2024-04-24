package grpc_hdwallet

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
	cfg    *config.MangerConfig
	// all GRPC handlers
	generateMnemonicHandlerSvc generateMnemonicHandlerService
	loadMnemonicHandlerSvc     loadMnemonicHandlerService
	unLoadMnemonicHandlerSvc   unLoadMnemonicHandlerService
	getDerivationAddressSvc    getDerivationAddressHandlerService
	getDerivationsAddressesSvc getDerivationsAddressesHandlerService
	signTransactionSvc         signTransactionHandlerService
}

func (h *grpcServerHandle) GenerateMnemonic(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	return h.generateMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadMnemonic(ctx context.Context,
	req *pbApi.LoadMnemonicRequest,
) (*pbApi.LoadMnemonicResponse, error) {
	return h.loadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) UnLoadMnemonic(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	return h.unLoadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddress(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	return h.getDerivationAddressSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddressByRange(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	return h.getDerivationsAddressesSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) SignTransaction(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	return h.signTransactionSvc.Handle(ctx, req)
}

// New instance of service
func New(ctx context.Context,
	loggerSrv *zap.Logger,
) pbApi.HdWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler").With(
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	//addrRespPool := &sync.Pool{New: func() any {
	//	return new(pbApi.DerivationAddressIdentity)
	//}}

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		generateMnemonicHandlerSvc: MakeGenerateMnemonicHandler(l),
		loadMnemonicHandlerSvc:     MakeLoadMnemonicHandler(l),
		unLoadMnemonicHandlerSvc:   MakeUnLoadMnemonicHandler(l),
		getDerivationAddressSvc:    MakeGetDerivationAddressHandler(l),
		getDerivationsAddressesSvc: MakeGetDerivationsAddressesHandler(l),
	}
}
