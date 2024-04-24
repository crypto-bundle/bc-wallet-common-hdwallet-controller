package grpc

import (
	"context"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"sync"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

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

	walletPubData, err := h.walletSvc.GetWalletByUUID(ctx, vf.MnemonicWalletUUIDRaw)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletPubData == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	addressData, err := h.walletSvc.GetAddressByPath(ctx, vf.MnemonicWalletUUIDRaw,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex)
	if err != nil {
		return nil, err
	}

	addressEntity := h.pbAddrPool.Get().(*pbCommon.DerivationAddressIdentity)
	addressEntity.AccountIndex = addressData.AccountIndex
	addressEntity.InternalIndex = addressData.InternalIndex
	addressEntity.AddressIndex = addressData.AddressIndex
	addressEntity.Address = addressData.Address

	marshalledData, err := h.marshallerSvc.MarshallGetAddressData(walletPubData, addressEntity)
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
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
	pbAddrPool *sync.Pool,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSvc:     walletSrv,
		marshallerSvc: marshallerSrv,
		pbAddrPool:    pbAddrPool,
	}
}
