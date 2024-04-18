package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"github.com/google/uuid"
)

type GetWalletSessionForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID

	SessionUUID    string `valid:"type(string),uuid,required"`
	SessionUUIDRaw uuid.UUID
}

func (f *GetWalletSessionForm) LoadAndValidate(ctx context.Context,
	req *pbApi.GetWalletSessionRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID

	if req.SessionIdentity == nil {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Session identity")
	}
	f.SessionUUID = req.SessionIdentity.SessionUUID

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	walletUUIDRaw, err := uuid.Parse(f.WalletUUID)
	if err != nil {
		return false, err
	}
	f.WalletUUIDRaw = walletUUIDRaw

	sessionUUIDRaw, err := uuid.Parse(f.SessionUUID)
	if err != nil {
		return false, err
	}
	f.SessionUUIDRaw = sessionUUIDRaw

	return true, nil
}
