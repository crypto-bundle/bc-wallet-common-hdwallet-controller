package handlers

import (
	"bc-wallet-eth-hdwallet/internal/forms"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bc-wallet-eth-hdwallet/internal/app"
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hd_wallet_api/proto"
	"github.com/crypto-bundle/bc-adapter-common/pkg/tracer"
)

const (
	MethodGetDerivationAddress = "GetDerivationAddress"
)

type GetDerivationAddressHandler struct {
	l         *zap.Logger
	walletSrv walleter
}

// nolint:funlen // fixme
func (h *GetDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	validationForm := &forms.GetDerivationAddressForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	address, err := h.walletSrv.GetAddressByPath(ctx, validationForm.WalletUUID,
		validationForm.AccountIndex, validationForm.InternalIndex, validationForm.AddressIndex)
	if err != nil {
		return nil, err
	}

	return &pbApi.DerivationAddressResponse{
		Address: address,
	}, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletSrv walleter,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSrv: walletSrv,
	}
}
