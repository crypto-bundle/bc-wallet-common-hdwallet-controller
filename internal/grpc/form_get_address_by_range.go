package grpc

import (
	"context"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

type derivationAddressByRangeForm struct {
	MnemonicWalletUUID    string `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID

	SessionUUID string `valid:"type(string),uuid,required"`

	Ranges      []*pbCommon.RangeRequestUnit `valid:"required"`
	RangesCount uint32                       `valid:"type(uint32),required"`
	RangeSize   uint32                       `valid:"type(uint32),required"`

	index uint32
}

func (f *derivationAddressByRangeForm) hasNext() bool {
	if f.index < f.RangesCount {
		return true
	}
	return false
}

func (f *derivationAddressByRangeForm) GetRangesCount() uint32 {
	return f.RangesCount
}

func (f *derivationAddressByRangeForm) GetRangesSize() uint32 {
	return f.RangeSize
}

func (f *derivationAddressByRangeForm) GetNext() *pbCommon.RangeRequestUnit {
	if f.hasNext() {
		rageForm := f.Ranges[f.index]
		f.index++

		return rageForm
	}
	return nil
}

func (f *derivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "MnemonicWallet identity")
	}
	f.MnemonicWalletUUID = req.MnemonicIdentity.WalletUUID

	if req.SessionIdentity == nil {
		return false, fmt.Errorf("%w:%s",
			ErrMissedRequiredData, "Session identity")
	}
	f.SessionUUID = req.SessionIdentity.SessionUUID

	if req.Ranges == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Ranges data")
	}
	f.RangesCount = uint32(len(req.Ranges))
	f.Ranges = req.Ranges

	for i := uint32(0); i != f.RangesCount; i++ {
		data := req.Ranges[i]
		diff := (data.AddressIndexTo - data.AddressIndexFrom) + 1
		if diff <= 0 {
			return false, fmt.Errorf("%w:%s, %d: %d", ErrDataIsNotValid,
				"Values diff in range lower or equal to 0", data.AddressIndexTo, data.AddressIndexFrom)
		}

		if data.AddressIndexTo == data.AddressIndexFrom {
			diff = 1
		}

		f.RangeSize += diff
	}

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
