package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameEnableWallets = "EnableWallets"
)

type EnableWalletsHandler struct {
	l *zap.Logger

	walletSvc walletManagerService
}

// nolint:funlen // fixme
func (h *EnableWalletsHandler) Handle(ctx context.Context,
	req *pbApi.EnableWalletsRequest,
) (*pbApi.EnableWalletsResponse, error) {
	var err error

	validationForm := &WalletsIdentitiesForm{}
	valid, err := validationForm.LoadAndValidate(req.WalletIdentities)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	disabledCount, walletsIdentities, err := h.walletSvc.EnableWalletsByUUIDList(ctx, validationForm.WalletUUIDs)
	if err != nil {
		h.l.Error("unable to disable wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if walletsIdentities == nil {
		return nil, status.Error(codes.NotFound, "there are no wallets available to enable")
	}

	pbIdentities := make([]*common.MnemonicWalletIdentity, disabledCount)
	for i := uint(0); i != disabledCount; i++ {
		pbIdentities[i] = &common.MnemonicWalletIdentity{
			WalletUUID: walletsIdentities[i],
		}
	}

	return &pbApi.EnableWalletsResponse{
		WalletIdentities: pbIdentities,
	}, nil
}

func MakeEnableWalletsHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *EnableWalletsHandler {
	return &EnableWalletsHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodNameEnableWallets)),
		walletSvc: walletSvc,
	}
}
