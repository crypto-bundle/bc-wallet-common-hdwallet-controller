package grpc

import (
	"context"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameImportWallet = "ImportWallet"
)

type ImportWalletHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *ImportWalletHandler) Handle(ctx context.Context,
	req *pbApi.ImportWalletRequest,
) (*pbApi.ImportWalletResponse, error) {
	var err error

	validationForm := &ImportWalletForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	wallet, err := h.walletSvc.ImportWallet(ctx, validationForm.Phrase)
	if err != nil {
		h.l.Error("unable to import mnemonic wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbApi.ImportWalletResponse{
		WalletIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: wallet.UUID.String(),
			WalletHash: wallet.MnemonicHash,
		},
	}, nil
}

func MakeImportWalletHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *ImportWalletHandler {
	return &ImportWalletHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameImportWallet)),

		walletSvc: walletSvc,
	}
}
