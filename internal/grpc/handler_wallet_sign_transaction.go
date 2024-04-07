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
	MethodSignTransaction = "SignTransaction"
)

type SignTransactionHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignTransactionHandler) Handle(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	var err error

	vf := &SignTransactionForm{}
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

	signOwner, signedTxData, err := h.walletSvc.SignTransaction(ctx, vf.WalletUUID,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex,
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

	return &pbApi.SignTransactionResponse{
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
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *SignTransactionHandler {
	return &SignTransactionHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodSignTransaction)),
		walletSvc:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
