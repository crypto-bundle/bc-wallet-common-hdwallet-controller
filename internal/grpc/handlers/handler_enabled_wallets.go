package handlers

import (
	"bc-wallet-eth-hdwallet/internal/app"
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	"context"
	"github.com/crypto-bundle/bc-adapter-common/pkg/tracer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetEnabledWallets = "GetEnabledWallets"
)

type GetEnabledWalletsHandler struct {
	l         *zap.Logger
	walletSrv walleter
}

// nolint:funlen // fixme
func (h *GetEnabledWalletsHandler) Handle(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	wallets, err := h.walletSrv.GetEnabledWalletsUUID(ctx)
	if err != nil {
		h.l.Error("unable to get enabled mnemonic wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbApi.GetEnabledWalletsResponse{
		WalletsUUID: wallets,
	}, nil
}

func MakeGetEnabledWalletsHandler(loggerEntry *zap.Logger,
	walletSrv walleter,
) *GetEnabledWalletsHandler {
	return &GetEnabledWalletsHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetEnabledWallets)),
		walletSrv: walletSrv,
	}
}
