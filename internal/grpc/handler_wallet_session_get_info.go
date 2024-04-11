package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetWalletSession = "GetWalletSession"
)

type GetWalletSessionHandler struct {
	l *zap.Logger

	walletSvc walletManagerService
}

// nolint:funlen // fixme
func (h *GetWalletSessionHandler) Handle(ctx context.Context,
	req *pbApi.GetWalletSessionRequest,
) (*pbApi.GetWalletSessionResponse, error) {
	var err error

	vf := &GetWalletSessionForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, sessionItem, err := h.walletSvc.GetWalletSessionInfo(ctx, vf.WalletUUID, vf.SessionUUID)
	if err != nil {
		h.l.Error("unable get wallet and wallet session info", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	if sessionItem == nil {
		return nil, status.Error(codes.ResourceExhausted, "mnemonic wallet session not found or expired")
	}

	return &pbApi.GetWalletSessionResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		Session: &pbApi.SessionInfo{
			SessionIdentity: &pbApi.WalletSessionIdentity{
				SessionUUID: sessionItem.UUID,
			},
			SessionStartedAt: uint64(sessionItem.CreatedAt.Unix()),
			SessionExpiredAt: uint64(sessionItem.ExpiredAt.Unix()),
		},
	}, nil
}

func MakeGetWalletSessionHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
) *GetWalletSessionHandler {
	return &GetWalletSessionHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodGetWalletSession)),

		walletSvc: walletSrv,
	}
}
