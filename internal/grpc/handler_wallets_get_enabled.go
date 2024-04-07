package grpc

import (
	"context"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

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
	_ *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	var err error

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
	walletSvc walletManagerService,
	marshallerSrv marshallerService,
) *GetEnabledWalletsHandler {
	return &GetEnabledWalletsHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetEnabledWallets)),
		walletSrv:     walletSvc,
		marshallerSrv: marshallerSrv,
	}
}
