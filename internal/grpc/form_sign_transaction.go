package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

type SignTransactionForm struct {
	WalletUUID  string `valid:"type(string),uuid,required"`
	SessionUUID string `valid:"type(string),uuid,required"`

	AccountIndex  uint32 `valid:"type(uint)"`
	InternalIndex uint32 `valid:"type(uint)"`
	AddressIndex  uint32 `valid:"type(uint)"`

	SignData []byte `valid:"type([]byte]),required"`
}

func (f *SignTransactionForm) LoadAndValidate(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID

	if req.SessionIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet sesssion identity")
	}
	f.SessionUUID = req.SessionIdentity.SessionUUID

	if req.AddressIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	f.AccountIndex = req.AddressIdentity.AccountIndex
	f.InternalIndex = req.AddressIdentity.InternalIndex
	f.AddressIndex = req.AddressIdentity.AddressIndex

	if req.CreatedTxData == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Missing signature data")
	}
	f.SignData = req.CreatedTxData

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
