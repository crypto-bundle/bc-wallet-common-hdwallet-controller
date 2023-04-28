package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type GetWalletInfoForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID
}

func (f *GetWalletInfoForm) LoadAndValidate(ctx context.Context,
	req *pbApi.GetWalletInfoRequest,
) (valid bool, err error) {
	if req.WalletIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
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
