package grpc

import (
	"context"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetWalletInfo = "GetWalletInfo"
)

type GetWalletInfoHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

func (h *GetWalletInfoHandler) Handle(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &GetWalletInfoForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletPubData, err := h.walletSrv.GetWalletByUUID(ctx, vf.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletPubData == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	walletInfo := h.marshallerSrv.MarshallWalletInfo(walletPubData)

	return &pbApi.GetWalletInfoResponse{
		WalletIdentity: walletInfo.Identity,
		WalletInfo:     walletInfo,
	}, nil
}

func MakeGetWalletInfoHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetWalletInfoHandler {
	return &GetWalletInfoHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetWalletInfo)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
