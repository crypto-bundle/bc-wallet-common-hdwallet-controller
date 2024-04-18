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
	MethodNameDisableWallet = "DisableWallet"
)

type DisableWalletHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService

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

	wallet, err := h.walletSvc.DisableWalletByUUID(ctx, validationForm.WalletUUID)
	if err != nil {
		h.l.Error("unable to disable wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	_, _, err = h.signManagerSvc.CloseSignRequestByWallet(ctx, wallet.UUID.String())
	if err != nil {
		h.l.Error("unable to close sign requests by session", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, wallet.UUID.String()))

		// no return err - it's ok
	}

	return &pbApi.DisableWalletResponse{
		WalletIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: wallet.UUID.String(),
			WalletHash: wallet.MnemonicHash,
		},
	}, nil
}

func MakeDisableWalletHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	signManagerSvc signManagerService,
) *DisableWalletHandler {
	return &DisableWalletHandler{
		l:              loggerEntry.With(zap.String(MethodNameTag, MethodNameDisableWallet)),
		walletSvc:      walletSvc,
		signManagerSvc: signManagerSvc,
	}
}
