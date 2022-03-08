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
	MethodGetDerivationAddressByRange = "GetDerivationAddressByRange"
)

type GetDerivationAddressByRangeHandler struct {
	l         *zap.Logger
	walletSrv walleter
}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	validationForm := &forms.DerivationAddressByRangeForm{}
	valid, err := validationForm.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	rangeSize := validationForm.AddressIndexTo - validationForm.AddressIndexFrom

	response := &pbApi.DerivationAddressByRangeResponse{
		AddressIdentities: make([]*pbApi.DerivationAddressIdentity, rangeSize+1),
	}

	for i, j := validationForm.AddressIndexFrom, uint32(0); i <= validationForm.AddressIndexTo; i++ {
		address, err := h.walletSrv.GetAddressByPath(ctx, validationForm.WalletUUID,
			validationForm.AccountIndex, validationForm.InternalIndex, i)
		if err != nil {
			return nil, err
		}

		response.AddressIdentities[j] = &pbApi.DerivationAddressIdentity{
			AccountIndex:  validationForm.AccountIndex,
			InternalIndex: validationForm.InternalIndex,
			AddressIndex:  i,
			Address:       address,
		}

		j++
	}

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSrv walleter,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSrv: walletSrv,
	}
}
