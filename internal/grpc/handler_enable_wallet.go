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
	MethodNameEnableWallet = "EnableWallet"
)

type EnableWalletHandler struct {
	l *zap.Logger

	walletDataSvc mnemonicWalletsDataService

	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *EnableWalletHandler) Handle(ctx context.Context,
	req *pbApi.EnableWalletRequest,
) (*pbApi.EnableWalletResponse, error) {
	var err error

	validationForm := &EnableWalletForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err),
			zap.String(app.WalletUUIDTag, req.WalletIdentity.WalletUUID))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	wallet, err := h.walletDataSvc.EnableWalletByUUID(ctx, validationForm.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable to enable mnemonic wallet", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbApi.EnableWalletResponse{
		WalletIdentity: &pbCommon.MnemonicWalletIdentity{WalletUUID: wallet.UUID.String()},
	}, nil
}

func MakeEnableWalletHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *AddNewWalletHandler {
	return &AddNewWalletHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameEnableWallet)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
