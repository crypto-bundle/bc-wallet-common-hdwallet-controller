package forms

import (
	"context"
	"errors"

	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hd_wallet_api/proto"
	protoTypes "bc-wallet-eth-hdwallet/pkg/types"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/metadata"
)

var (
	ErrUnableReadGrpcMetadata          = errors.New("unable to read grpc metadata")
	ErrUnableGetWalletUUIDFromMetadata = errors.New("unable to get wallet uuid from metadata")
)

type GetDerivationAddressForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	AccountIndex  uint32 `valid:"type(uint32),int,required"`
	InternalIndex uint32 `valid:"type(uint32),int,required"`
	AddressIndex  uint32 `valid:"type(uint32),int,required"`
}

func (f *GetDerivationAddressForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
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
	f.AddressIndex = req.AccountIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
