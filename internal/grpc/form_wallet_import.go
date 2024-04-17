package grpc

import (
	"context"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

type ImportWalletForm struct {
	Phrase []byte `valid:"required"`
}

func (f *ImportWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.ImportWalletRequest,
) (valid bool, err error) {
	f.Phrase = req.MnemonicPhrase

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
