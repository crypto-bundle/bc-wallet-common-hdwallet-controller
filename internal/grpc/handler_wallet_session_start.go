package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"sync"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodStartWalletSession = "StartWalletSession"
)

type StartWalletSessionHandler struct {
	l *zap.Logger

	walletSvc     walletManagerService
	marshallerSvc marshallerService

	pbAddrPool *sync.Pool
}

// nolint:funlen // fixme
func (h *StartWalletSessionHandler) Handle(ctx context.Context,
	req *pbApi.StartWalletSessionRequest,
) (*pbApi.StartWalletSessionResponse, error) {
	var err error

	vf := &StartWalletSessionForm{}
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

	if walletItem.Status == types.MnemonicWalletStatusDisabled {
		return nil, status.Error(codes.ResourceExhausted, "wallet disabled")
	}

	sessionItem, err := h.walletSvc.StartSessionForWallet(ctx, walletItem)
	if err != nil {
		h.l.Error("unable to start wallet session", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if sessionItem == nil || !sessionItem.IsSessionActive() {
		return nil, status.Error(codes.ResourceExhausted, "wallet session not found or already expired")
	}

	return &pbApi.StartWalletSessionResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentity: &pbApi.WalletSessionIdentity{
			SessionUUID: sessionItem.UUID,
		},
		SessionStartedAt: uint64(sessionItem.CreatedAt.Unix()),
		SessionExpiredAt: uint64(sessionItem.ExpiredAt.Unix()),
	}, nil
}

func MakeStartWalletSessionHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *StartWalletSessionHandler {
	return &StartWalletSessionHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodStartWalletSession)),
		walletSvc: walletSvc,
	}
}
