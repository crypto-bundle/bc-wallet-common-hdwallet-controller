package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
	"github.com/google/uuid"
)

type CloseWalletSessionForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID

	SessionUUID    string `valid:"type(string),uuid,required"`
	SessionUUIDRaw uuid.UUID
}

func (f *CloseWalletSessionForm) LoadAndValidate(ctx context.Context,
	req *pbApi.CloseWalletSessionsRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData,
			"Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID

	if req.SessionIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData,
			"Session identity")
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

	return true, nil
}
