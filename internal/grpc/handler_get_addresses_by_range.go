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
	MethodGetDerivationAddressByRange = "GetDerivationAddressByRange"
)

type GetDerivationAddressByRangeHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error

	vf := &derivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return h.processRequest(ctx, vf)
}

func (h *GetDerivationAddressByRangeHandler) processRequest(ctx context.Context,
	vf *derivationAddressByRangeForm,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error

	owner, addressesList, err := h.walletSrv.GetAddressesByRange(ctx, vf.MnemonicWalletUUID, vf.SessionUUID,
		vf.Ranges)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	response, err := h.marshallerSrv.MarshallGetAddressByRange(owner, addressesList,
		uint64(vf.RangeSize))
	if err != nil {
		h.l.Error("unable to marshall get addresses data", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.MnemonicWalletUUID),
			zap.String(app.MnemonicWalletSessionUUIDTag, vf.SessionUUID))

		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSvc marshallerService,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSrv:     walletSvc,
		marshallerSrv: marshallerSvc,
	}
}
