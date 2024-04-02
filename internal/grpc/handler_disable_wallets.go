package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameDisableWallets = "DisableWallets"
)

type DisableWalletsHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
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

	disabledCount, walletsIdentities, err := h.walletSrv.DisableWalletsByUUIDList(ctx, validationForm.WalletUUIDs)
	if err != nil {
		h.l.Error("unable to disable wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if walletsIdentities == nil {
		return nil, status.Error(codes.NotFound, "there are no wallets available to disable")
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
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *DisableWalletHandler {
	return &DisableWalletHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameDisableWallets)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
