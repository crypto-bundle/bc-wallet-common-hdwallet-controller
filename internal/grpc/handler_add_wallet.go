package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
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
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *AddNewWalletHandler) Handle(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	var err error

	validationForm := &AddNewWalletForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err),
			zap.String(app.WalletTitleTag, req.Title),
			zap.String(app.WalletPurposeTag, req.Purpose))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	wallet, err := h.walletSrv.CreateNewWallet(ctx,
		validationForm.Title, validationForm.Purpose)
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
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *AddNewWalletHandler {
	return &AddNewWalletHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameAddNewWallet)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
