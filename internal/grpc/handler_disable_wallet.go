package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameDisableWallet = "DisableWallet"
)

type DisableWalletHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *DisableWalletHandler) Handle(ctx context.Context,
	req *pbApi.DisableWalletRequest,
) (*pbApi.DisableWalletResponse, error) {
	var err error

	validationForm := &DisableWalletForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err),
			zap.String(app.WalletUUIDTag, req.WalletIdentity.WalletUUID))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	wallet, err := h.walletSrv.DisableWalletByUUID(ctx, validationForm.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable to disable wallet", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbApi.DisableWalletResponse{
		WalletIdentity: &pbCommon.MnemonicWalletIdentity{WalletUUID: wallet.UUID.String()},
	}, nil
}

func MakeDisableWalletHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *DisableWalletHandler {
	return &DisableWalletHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameDisableWallet)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
