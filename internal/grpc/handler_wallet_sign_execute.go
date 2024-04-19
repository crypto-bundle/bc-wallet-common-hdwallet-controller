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
	MethodSignTransaction = "SignTransaction"
)

type SignTransactionHandler struct {
	l *zap.Logger

	walletSvc      walletManagerService
	signManagerSvc signManagerService

	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignTransactionHandler) Handle(ctx context.Context,
	req *pbApi.ExecuteSignRequestReq,
) (*pbApi.ExecuteSignRequestResponse, error) {
	var err error

	vf := &SignRequestExecForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	signReqItem, err := h.signManagerSvc.GetActiveSignRequest(ctx, vf.SignRequestUUID)
	if err != nil {
		h.l.Error("unable to get sign request info", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID),
			zap.String(app.SignRequestUUIDTag, vf.SignRequestUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if signReqItem == nil {
		return nil, status.Error(codes.NotFound, "sign request not found or already processed")
	}

	if vf.WalletUUID != signReqItem.WalletUUID {
		return nil, status.Error(codes.InvalidArgument, "mismatched wallet uuid")
	}

	if vf.SessionUUID != signReqItem.SessionUUID {
		return nil, status.Error(codes.InvalidArgument, "mismatched session uuid")
	}

	walletItem, sessionItem, err := h.walletSvc.GetWalletSessionInfo(ctx, signReqItem.WalletUUID,
		signReqItem.SessionUUID)
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

	signOwner, signedTxData, err := h.signManagerSvc.ExecuteSignRequest(ctx, signReqItem,
		vf.SignData)
	if err != nil {
		h.l.Error("unable to sign transaction", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, err.Error())
	}

	if signOwner == nil || signedTxData == nil {
		return nil, status.Error(codes.ResourceExhausted,
			"signer account not found or signature session expired")
	}

	return &pbApi.ExecuteSignRequestResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletItem.UUID.String(),
			WalletHash: walletItem.MnemonicHash,
		},
		SessionIdentity: &pbApi.WalletSessionIdentity{
			SessionUUID: sessionItem.UUID,
		},
		TxOwnerIdentity: signOwner,
		SignedTxData:    signedTxData,
	}, nil
}

func MakeSignTransactionsHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	signManagerSvc signManagerService,
	marshallerSrv marshallerService,
) *SignTransactionHandler {
	return &SignTransactionHandler{
		l:              loggerEntry.With(zap.String(MethodNameTag, MethodSignTransaction)),
		walletSvc:      walletSvc,
		signManagerSvc: signManagerSvc,
		marshallerSrv:  marshallerSrv,
	}
}
