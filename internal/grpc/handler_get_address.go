package grpc

import (
	"context"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"sync"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddress = "GetDerivationAddress"
)

type GetDerivationAddressHandler struct {
	l *zap.Logger

	walletSvc     walletManagerService
	marshallerSvc marshallerService

	pbAddrPool *sync.Pool
}

// nolint:funlen // fixme
func (h *GetDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error

	vf := &GetDerivationAddressForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	ownerWallet, addressData, err := h.walletSvc.GetAddress(ctx, vf.MnemonicWalletUUID,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex, req.SessionIdentity.SessionUUID)
	if err != nil {
		return nil, err
	}

	if ownerWallet == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	addressEntity := h.pbAddrPool.Get().(*pbCommon.DerivationAddressIdentity)
	addressEntity.AccountIndex = vf.AccountIndex
	addressEntity.InternalIndex = vf.InternalIndex
	addressEntity.AddressIndex = vf.AddressIndex
	addressEntity.Address = *addressData

	marshalledData, err := h.marshallerSvc.MarshallGetAddressData(ownerWallet, addressEntity)
	if err != nil {
		h.l.Error("unable to marshall public address data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer func() {
		h.pbAddrPool.Put(addressEntity)
	}()

	return marshalledData, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
	marshallerSrv marshallerService,
	pbAddrPool *sync.Pool,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSvc:     walletSvc,
		marshallerSvc: marshallerSrv,
		pbAddrPool:    pbAddrPool,
	}
}
