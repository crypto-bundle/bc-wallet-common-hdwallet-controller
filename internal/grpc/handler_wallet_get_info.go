package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetWalletInfo = "GetWalletInfo"
)

type GetWalletInfoHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

func (h *GetWalletInfoHandler) Handle(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (*pbApi.GetWalletInfoResponse, error) {
	var err error

	vf := &GetWalletInfoForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, err := h.walletSvc.GetWalletByUUID(ctx, vf.WalletUUID)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	return &pbApi.GetWalletInfoResponse{
		WalletIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
	}, nil
}

func MakeGetWalletInfoHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *GetWalletInfoHandler {
	return &GetWalletInfoHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodGetWalletInfo)),

		walletSvc: walletSvc,
	}
}
