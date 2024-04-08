package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

type SignRequestPrepareForm struct {
	WalletUUID  string `valid:"type(string),uuid,required"`
	SessionUUID string `valid:"type(string),uuid,required"`
	PurposeUUID string `valid:"type(string),uuid,required"`

	AccountIndex  uint32 `valid:"type(uint)"`
	InternalIndex uint32 `valid:"type(uint)"`
	AddressIndex  uint32 `valid:"type(uint)"`
}

func (f *SignRequestPrepareForm) LoadAndValidate(ctx context.Context,
	req *pbApi.PrepareSignRequestReq,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID

	if req.SessionIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet sesssion identity")
	}
	f.SessionUUID = req.SessionIdentity.SessionUUID

	if req.SignPurposeIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Signature puprose identity")
	}
	f.PurposeUUID = req.SignPurposeIdentifier.UUID

	if req.AddressIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	f.AccountIndex = req.AddressIdentity.AccountIndex
	f.InternalIndex = req.AddressIdentity.InternalIndex
	f.AddressIndex = req.AddressIdentity.AddressIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
