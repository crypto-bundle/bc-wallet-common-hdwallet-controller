package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetWalletSessions = "GetWalletSessions"
)

type GetWalletSessionsHandler struct {
	l *zap.Logger

	walletSvc     walletManagerService
	marshallerSvc marshallerService
}

// nolint:funlen // fixme
func (h *GetWalletSessionsHandler) Handle(ctx context.Context,
	req *pbApi.GetWalletSessionsRequest,
) (*pbApi.GetWalletSessionsResponse, error) {
	var err error

	vf := &GetWalletSessionsForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, sessionsList, err := h.walletSvc.GetWalletSessionsByWalletUUID(ctx, vf.WalletUUID)
	if err != nil {
		h.l.Error("unable to get wallet and all wallets sessions", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	if walletItem.Status == types.MnemonicWalletStatusDisabled {
		return nil, status.Error(codes.ResourceExhausted, "wallet disabled")
	}

	if sessionsList == nil {
		return nil, status.Error(codes.ResourceExhausted, "active wallet sessions not found or already expired")
	}

	return &pbApi.GetWalletSessionsResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		ActiveSessions: h.marshallerSvc.MarshallWalletSessions(sessionsList),
	}, nil
}

func MakeGetWalletSessionsHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetWalletSessionsHandler {
	return &GetWalletSessionsHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetWalletSessions)),
		walletSvc:     walletSrv,
		marshallerSvc: marshallerSrv,
	}
}
