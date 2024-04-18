package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/asaskevich/govalidator"
)

type DisableWalletsForm struct {
	WalletUUIDs []string `valid:"type([]string),required"`
}

func (f *DisableWalletsForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DisableWalletsRequest,
) (valid bool, err error) {
	if len(req.WalletIdentities) == 0 {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identities")
	}

	f.WalletUUIDs = make([]string, len(req.WalletIdentities))

	for i, v := range req.WalletIdentities {
		if !govalidator.IsUUID(v.WalletUUID) {
			return false, fmt.Errorf("%s does not validate as %s", v.WalletUUID, "UUID")
		}

		f.WalletUUIDs[i] = v.WalletUUID
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
