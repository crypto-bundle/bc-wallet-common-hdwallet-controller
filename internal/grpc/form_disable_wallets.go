package grpc

import (
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	"github.com/asaskevich/govalidator"
)

type WalletsIdentitiesForm struct {
	WalletUUIDs []string `valid:"type([]string),required"`
}

func (f *WalletsIdentitiesForm) LoadAndValidate(
	list []*pbCommon.MnemonicWalletIdentity,
) (valid bool, err error) {
	count := len(list)
	if count == 0 {
		return false,
			fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identities")
	}

	f.WalletUUIDs = make([]string, count)

	for i, v := range list {
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
