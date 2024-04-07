package grpc

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameAddNewWallet = "AddNewWallet"
)

type AddNewWalletHandler struct {
	l             *zap.Logger
	walletSvc     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *AddNewWalletHandler) Handle(ctx context.Context,
	_ *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	var err error

	wallet, err := h.walletSvc.AddNewWallet(ctx)
	if err != nil {
		h.l.Error("unable to create mnemonic wallet", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	marshalledData, err := h.marshallerSrv.MarshallCreateWalletData(wallet)
	if err != nil {
		h.l.Error("unable to marshall wallet data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return marshalledData, nil
}

func MakeAddNewWalletHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSrv marshallerService,
) *AddNewWalletHandler {
	return &AddNewWalletHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameAddNewWallet)),
		walletSvc:     walletSvc,
		marshallerSrv: marshallerSrv,
	}
}
