package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/manager"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type EnableWalletForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID
}

func (f *EnableWalletForm) LoadAndValidate(ctx context.Context,
	req *pbApi.EnableWalletRequest,
) (valid bool, err error) {
	if req.WalletIdentity == nil {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.WalletIdentity.WalletUUID

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	walletUUIDRaw, err := uuid.Parse(f.WalletUUID)
	if err != nil {
		return false, err
	}
	f.WalletUUIDRaw = walletUUIDRaw

	return true, nil
}
