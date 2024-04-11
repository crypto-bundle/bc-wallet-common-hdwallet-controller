package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameDisableWallets = "DisableWallets"
)

type DisableWalletsHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService
}

// nolint:funlen // fixme
func (h *DisableWalletsHandler) Handle(ctx context.Context,
	req *pbApi.DisableWalletsRequest,
) (*pbApi.DisableWalletsResponse, error) {
	var err error

	validationForm := &DisableWalletsForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	disabledCount, walletsIdentities, err := h.walletSvc.DisableWalletsByUUIDList(ctx, validationForm.WalletUUIDs)
	if err != nil {
		h.l.Error("unable to disable wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if walletsIdentities == nil {
		return nil, status.Error(codes.NotFound, "there are no wallets available to disable")
	}

	_, _, err = h.signManagerSvc.CloseSignRequestByMultipleWallets(ctx, walletsIdentities)
	if err != nil {
		h.l.Error("unable to close sign requests by session", zap.Error(err),
			zap.Strings(app.MnemonicWalletUUIDTag, walletsIdentities))

		// no return err - it's ok
	}

	pbIdentities := make([]*common.MnemonicWalletIdentity, disabledCount)
	for i := uint(0); i != disabledCount; i++ {
		pbIdentities[i] = &common.MnemonicWalletIdentity{
			WalletUUID: walletsIdentities[i],
		}
	}

	return &pbApi.DisableWalletsResponse{
		WalletIdentities: pbIdentities,
	}, nil
}

func MakeDisableWalletsHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	signManagerSvc signManagerService,
) *DisableWalletsHandler {
	return &DisableWalletsHandler{
		l:              loggerEntry.With(zap.String(MethodNameTag, MethodNameDisableWallets)),
		walletSvc:      walletSvc,
		signManagerSvc: signManagerSvc,
	}
}
