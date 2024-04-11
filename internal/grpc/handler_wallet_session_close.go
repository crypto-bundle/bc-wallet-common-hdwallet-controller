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
	MethodCloseWalletSession = "CloseWalletSession"
)

type CloseWalletSessionHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService
	marshallerSvc  marshallerService
}

// nolint:funlen // fixme
func (h *CloseWalletSessionHandler) Handle(ctx context.Context,
	req *pbApi.CloseWalletSessionsRequest,
) (*pbApi.CloseWalletSessionsResponse, error) {
	var err error

	vf := &CloseWalletSessionForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletItem, sessionItem, err := h.walletSvc.CloseWalletSession(ctx, vf.WalletUUID, vf.SessionUUID)
	if err != nil {
		h.l.Error("unable to start wallet session", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletItem == nil || sessionItem == nil {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	_, err = h.signManagerSvc.CloseSignRequestBySession(ctx, vf.SessionUUID)
	if err != nil {
		h.l.Error("unable to close sign requests by session", zap.Error(err),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		// no return err - it's ok
	}

	return &pbApi.CloseWalletSessionsResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentity: &pbApi.WalletSessionIdentity{
			SessionUUID: sessionItem.UUID,
		},
	}, nil
}

func MakeCloseWalletSessionHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	signManagerSvc signManagerService,
) *CloseWalletSessionHandler {
	return &CloseWalletSessionHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodCloseWalletSession)),

		walletSvc:      walletSvc,
		signManagerSvc: signManagerSvc,
	}
}
