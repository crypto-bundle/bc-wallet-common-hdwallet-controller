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
	MethodSignPrepare = "PrepareSign"
)

type SignPrepareHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignPrepareHandler) Handle(ctx context.Context,
	req *pbApi.PrepareSignRequest,
) (*pbApi.PrepareSignResponse, error) {
	var err error

	vf := &SignPrepareForm{}
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

	signOwner, signReqItem, err := h.walletSvc.PrepareForSign(ctx, vf.WalletUUID,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex)
	if err != nil {
		h.l.Error("unable to sign transaction", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, err.Error())
	}

	if signOwner == nil || signReqItem == nil {
		return nil, status.Error(codes.ResourceExhausted,
			"signer account not found or signature session expired")
	}

	return &pbApi.PrepareSignResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentity: &pbApi.WalletSessionIdentity{
			SessionUUID: sessionItem.UUID,
		},
		TxOwnerIdentity: signOwner,
		SignatureRequestInfo: &pbApi.SignRequestData{
			Identifier: &pbApi.SignRequestIdentity{UUID: signReqItem.UUID},
			Status:     pbApi.SignRequestData_ReqStatus(signReqItem.Status),
			CreateAt:   uint64(signReqItem.CreatedAt.Unix()),
		},
	}, nil
}

func MakeSignPrepareHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *SignPrepareHandler {
	return &SignPrepareHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodSignPrepare)),
		walletSvc:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
