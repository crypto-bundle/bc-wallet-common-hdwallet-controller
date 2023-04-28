package grpc

import (
	"context"
	"errors"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"github.com/asaskevich/govalidator"
)

var (
	ErrUnableReadCreateWalletStrategy = errors.New("unable to reade wallet strategy")
)

type AddNewWalletForm struct {
	Title    string                    `valid:"type(string),required"`
	Purpose  string                    `valid:"type(string),required"`
	Strategy types.WalletMakerStrategy `valid:"type(types.WalletMakerStrategy)"`
}

func (f *AddNewWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.AddNewWalletRequest,
) (valid bool, err error) {
	f.Title = req.Title
	f.Purpose = req.Purpose
	f.Strategy = types.WalletMakerStrategy(req.Strategy)

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
