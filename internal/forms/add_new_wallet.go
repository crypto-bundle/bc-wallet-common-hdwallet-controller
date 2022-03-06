package forms

import (
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hd_wallet_api/proto"
	"context"

	"github.com/asaskevich/govalidator"
)

type AddNewWalletForm struct {
	Title   string `valid:"type(string),required"`
	Purpose string `valid:"type(string),required"`
	IsHot   bool   `valid:"type(bool),required"`
}

func (f *AddNewWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (valid bool, err error) {
	f.Title = req.Title
	f.Purpose = req.Purpose
	f.IsHot = req.IsHot

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
