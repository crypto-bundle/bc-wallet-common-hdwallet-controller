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
	MethodSignPrepare = "PrepareSign"
)

type SignPrepareHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService

	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignPrepareHandler) Handle(ctx context.Context,
	req *pbApi.PrepareSignRequestReq,
) (*pbApi.PrepareSignRequestResponse, error) {
	var err error

	vf := &SignRequestPrepareForm{}
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

	if !sessionItem.IsSessionActive() {
		return nil, status.Error(codes.ResourceExhausted, "mnemonic wallet session not found or expired")
	}

	signOwner, signReqItem, err := h.signManagerSvc.PrepareSignRequest(ctx, vf.WalletUUID, vf.SessionUUID,
		vf.PurposeUUID,
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

	if signReqItem.DerivationPath == nil {

	}

	return &pbApi.PrepareSignRequestResponse{
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
	signManagerSvc signManagerService,
	marshallerSrv marshallerService,
) *SignPrepareHandler {
	return &SignPrepareHandler{
		l:              loggerEntry.With(zap.String(MethodNameTag, MethodSignPrepare)),
		signManagerSvc: signManagerSvc,
		walletSvc:      walletSrv,
		marshallerSrv:  marshallerSrv,
	}
}
