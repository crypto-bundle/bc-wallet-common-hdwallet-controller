package forms

import (
	"context"
	"errors"

	pbApi "github.com/crypto-bundle/bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	protoTypes "github.com/crypto-bundle/bc-wallet-eth-hdwallet/pkg/types"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/metadata"
)

var (
	ErrUnableReadGrpcMetadata          = errors.New("unable to read grpc metadata")
	ErrUnableGetWalletUUIDFromMetadata = errors.New("unable to get wallet uuid from metadata")
)

type GetDerivationAddressForm struct {
	WalletUUID    string `valid:"type(string),uuid,required"`
	AccountIndex  uint32 `valid:"type(uint32)"`
	InternalIndex uint32 `valid:"type(uint32)"`
	AddressIndex  uint32 `valid:"type(uint32)"`
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
	f.AccountIndex = req.AddressIdentity.AccountIndex
	f.InternalIndex = req.AddressIdentity.InternalIndex
	f.AddressIndex = req.AddressIdentity.AccountIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
