package grpc

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"

	"github.com/asaskevich/govalidator"
)

type AddNewWalletForm struct {
	Title   string `valid:"type(string),required"`
	Purpose string `valid:"type(string),required"`
}

func (f *AddNewWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (valid bool, err error) {
	f.Title = req.Title
	f.Purpose = req.Purpose

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
