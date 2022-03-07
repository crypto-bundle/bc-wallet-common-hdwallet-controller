package forms

import (
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	protoTypes "bc-wallet-eth-hdwallet/pkg/types"
	"context"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/metadata"
)

type DerivationAddressByRangeForm struct {
	WalletUUID string `valid:"type(string),uuid,required"`

	AccountIndex     uint32 `valid:"type(uint32)"`
	InternalIndex    uint32 `valid:"type(uint32)"`
	AddressIndexFrom uint32 `valid:"type(uint32)"`
	AddressIndexTo   uint32 `valid:"type(uint32)"`
}

func (f *DerivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, ErrUnableReadGrpcMetadata
	}

	walletHeaders := headers.Get(protoTypes.WalletUUIDHeaderName)

	if len(walletHeaders) == 0 {
		return false, ErrUnableGetWalletUUIDFromMetadata
	}

	f.WalletUUID = walletHeaders[0]

	f.AccountIndex = req.AccountIndex
	f.InternalIndex = req.InternalIndex
	f.AddressIndexFrom = req.AddressIndexFrom
	f.AddressIndexTo = req.AddressIndexTo

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
