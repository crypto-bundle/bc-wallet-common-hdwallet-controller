package handlers

import (
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/forms"
	pbApi "github.com/crypto-bundle/bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	"context"
	"github.com/crypto-bundle/bc-adapter-common/pkg/tracer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameAddNewWallet = "AddNewWallet"
)

type AddNewWalletHandler struct {
	l         *zap.Logger
	walletSrv walleter
}

// nolint:funlen // fixme
func (h *AddNewWalletHandler) Handle(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (*pbApi.AddNewWalletResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	validationForm := &forms.AddNewWalletForm{}
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

	wallet, err := h.walletSrv.CreateNewMnemonicWallet(ctx, req.Title, req.Purpose, req.IsHot)
	if err != nil {
		h.l.Error("unable to create mnemonic wallet", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbApi.AddNewWalletResponse{
		WalletUUID: wallet.UUID.String(),
	}, nil
}

func MakeAddNewWalletHandler(loggerEntry *zap.Logger,
	walletSrv walleter,
) *AddNewWalletHandler {
	return &AddNewWalletHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodNameAddNewWallet)),
		walletSrv: walletSrv,
	}
}
