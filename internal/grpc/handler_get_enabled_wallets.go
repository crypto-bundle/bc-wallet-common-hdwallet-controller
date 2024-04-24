package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetEnabledWallets = "GetEnabledWallets"
)

type GetEnabledWalletsHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *GetEnabledWalletsHandler) Handle(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	wallets, err := h.walletSrv.GetEnabledWallets(ctx)
	if err != nil {
		h.l.Error("unable to get enabled mnemonic wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if wallets == nil {
		return nil, status.Error(codes.NotFound, "hdwallet-service has no active wallets")
	}

	response, err := h.marshallerSrv.MarshallGetEnabledWallets(wallets)
	if err != nil {
		h.l.Error("unable to marshall get wallets data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func MakeGetEnabledWalletsHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetEnabledWalletsHandler {
	return &GetEnabledWalletsHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetEnabledWallets)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
